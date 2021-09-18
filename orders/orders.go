package orders

import (
	"distribution-bridge/env"
	"distribution-bridge/global"
	"distribution-bridge/http"
	"distribution-bridge/logger"
	"distribution-bridge/products"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type Job struct {
	ID string
	Since *time.Time
	ProductsJob products.Job
	RequestManager *global.RequestManager
}

// SyncOrders runs through new orders and orders updates
// since :: It will process orders after this moment. If nil, then all orders.
func (j *Job) SyncOrders() {
	// Get new orders from seller account (Retailer side)
	j.SyncNewOrders()

	// Get order updates from buyer account (Supplier side)
	j.SyncOrderUpdates()
}

func (j *Job) SyncNewOrders() {
	logger.Info(j.ID, global.DomainOrders, "Syncing new orders...")

	page := 0
	allOrdersFound := false
	ordersCount := 0
	for !allOrdersFound {
		orders, err := j.getSellerNonShippedOrders(page)
		if err != nil {
			logger.Error(j.ID, global.DomainOrders, fmt.Sprintf("failed to get orders on page %d", page), err)
			return
		}

		ordersCount = ordersCount + len(orders)
		if len(orders) == 0 {
			logger.Info(j.ID, global.DomainOrders, fmt.Sprintf("All new orders have been found and synced [%d]", ordersCount))
			allOrdersFound = true
			continue
		}
		logger.Info(j.ID, global.DomainOrders, fmt.Sprintf(fmt.Sprintf("%d of %d orders been process", len(orders), ordersCount)))

		for _, order := range orders {
			j.ProcessNewOrder(order)
		}

		page++
	}
}

func (j *Job) ProcessNewOrder(order Order) {
	// Check if exist on buyer/supplier side using the seller order code against the buyer order code
	_, exists, err := j.getBuyerOrderWithBuyerOrderCode(order.SellerOrderCode)
	if err != nil {
		logger.Error(j.ID, global.DomainOrders, "failed to get order with buyer order code", err)
		return
	}

	if exists {
		logger.Info(j.ID, global.DomainOrders, fmt.Sprintf("Order has already be created :: %s", order.ID))
		return
	}
	// Create new instance of the order on the buyer side
	buyerOrder, err := j.ConvertToBuyerOrder(order)
	if err != nil {
		logger.Error(j.ID, global.DomainOrders, fmt.Sprintf("Failed to convert order to buyer order for %s (Seller Order ID)", order.ID), err)
		return
	}
	buyerOrderID, err := j.postNewBuyerOrderToAPI(buyerOrder)
	if err != nil {
		logger.Error(j.ID, global.DomainOrders, fmt.Sprintf("Failed to create new order for %s (Seller Order ID)", order.ID), err)
		return
	}
	logger.Info(j.ID, global.DomainOrders, fmt.Sprintf("New order created on the buyer account :: %s --> %s", order.ID, buyerOrderID))
}

func (j *Job) SyncOrderUpdates() {
	logger.Info(j.ID, global.DomainOrders, "Syncing new order updates...")
	page := 0
	allOrdersFound := false
	ordersCount := 0

	for !allOrdersFound {
		// Retrieve orders from buyer account
		buyerOrders, err := j.getBuyerShippedOrders(page)
		if err != nil {
			logger.Error(j.ID, global.DomainOrders, fmt.Sprintf("failed to get orders on page :: %d", page), err)
			return
		}

		ordersCount = ordersCount + len(buyerOrders)
		if len(buyerOrders) == 0 {
			logger.Info(j.ID, global.DomainOrders, fmt.Sprintf("All orders have been found [%d]", ordersCount))
			allOrdersFound = true
			continue
		}
		logger.Info(j.ID, global.DomainOrders, fmt.Sprintf("%d of %d orders that have shipped from the buyer", len(buyerOrders), ordersCount))

		for _, buyerOrder := range buyerOrders {
			j.ProcessOrderUpdate(buyerOrder)
		}

		page++
	}
}

func (j *Job) ProcessOrderUpdate(buyerOrder Order) {
	// Fetch the order
	order, exists, err := j.getSellerOrderWithSellerOrderCode(buyerOrder.BuyerOrderCode)
	if err != nil {
		logger.Error(j.ID, global.DomainOrders, "failed to get order with buyer order code", err)
		return
	}

	if !exists {
		logger.Error(j.ID, global.DomainOrders, "Order has not been synced to seller account", errors.New("error: order missing"))
		return
	}

	// Check if the order has been marked as shipped on the seller account (retail side)
	//if hasBuyerOrderShipped(buyerOrder) && !order.Shipped {
	if buyerOrder.Shipped && !order.Shipped {
		logger.Info(j.ID, global.DomainOrders, "Order has been shipped in buyer account, sharing it with the seller account")

		err := j.createFulfillmentOnSellerOrder(order.ID, buyerOrder.Fulfillments)
		if err != nil {
			logger.Error(j.ID, global.DomainOrders, "failed to create fulfillment on the seller order", err)
			return
		}
		logger.Info(j.ID, global.DomainOrders, "Order has been marked as shipped in both accounts")
	} else if !buyerOrder.Shipped && order.Shipped {
		logger.Error(j.ID, global.DomainOrders, "Order was marked as shipped in seller account but not buyer account", errors.New("invalid state"))
		return
	}
}

func (j *Job) createFulfillmentOnSellerOrder(orderID string, fulfillments []Fulfillment) error {
	for index, fulfillment := range fulfillments {
		newFulfillmentItems := []NewFulfillmentItem{}
		for _, newFulfillmentItem := range fulfillment.Items {
			newFulfillmentItems = append(newFulfillmentItems, NewFulfillmentItem{
				ID: int32(index + 1),
				SKU: newFulfillmentItem.Sku,
				Quantity: int32(newFulfillmentItem.Quantity),
			})
		}
		jsonPayload, err := json.Marshal(NewFulfillmentRequestBody{
			Carrier: fulfillment.Carrier,
			TrackingCode: fulfillment.TrackingCode,
			TrackingURLs: fulfillment.TrackingUrls,
			Items: newFulfillmentItems,
		})
		if err != nil {
			return err
		}
		fmt.Printf("jsonPayload :: %+v\n", string(jsonPayload)) // TODO KEITH

		_, err = http.PostRequest(fmt.Sprintf("/orders/%s/fulfillments", orderID), env.GetSellerAPIKey(), jsonPayload)
		if err != nil {
			return err
		}
	}
	return nil
}

// hasBuyerOrderShipped :: is a helper method for checking if all "seller orders" attached to a buyer order has shipped.
func hasBuyerOrderShipped(buyerOrder BuyerOrder) bool {
	if len(buyerOrder.SellerOrders) == 0 && len(buyerOrder.Items) > 0 {
		// There is a least one item on the order that needs to be shipped
		return false
	}
	for _, sellerOrder := range buyerOrder.SellerOrders {
		if !sellerOrder.Fulfilled {
			// Item(s) has not been fulfilled
			return false
		}
	}
	return true
}

// ConvertToBuyerOrder :: Converts an order from the seller order model to the buyer order model
func (j *Job) ConvertToBuyerOrder(o Order) (BuyerOrder, error) {
	buyerItems := []BuyerItem{}
	for _, item := range o.Items {
		// Look up the ID of the variant
		idOfVariant, err := j.ProductsJob.GetIDOfVariantBySellerBarcode(env.GetBuyerAPIKey(), item.Barcode, item.BarcodeType)
		if err != nil {
			return BuyerOrder{}, err
		}
		buyerItems = append(buyerItems, BuyerItem{
			VariantID: idOfVariant,
			BuyerReference: item.ID,
			Quantity: item.Quantity,
		})
	}
	return BuyerOrder{
		BuyerReference: o.SellerOrderCode,
		OrderedDate: o.Created,
		Created: o.Created,
		Updated: o.Updated,
		Address: o.ShippingAddress, // Does not support both billing and shipping
		Items: buyerItems,
	}, nil
}