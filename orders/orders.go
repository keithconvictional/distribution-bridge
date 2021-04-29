package orders

import (
	"distribution-bridge/env"
	"distribution-bridge/http"
	"distribution-bridge/logger"
	"distribution-bridge/products"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// Sync new orders from buyer account to seller account. Sync order updates both ways.
func SyncOrders() {
	// Get new orders from seller account (Retailer side)
	//syncNewOrders()

	// Get order updates from buyer account (Supplier side)
	syncOrderUpdates()
}

func syncOrderUpdates() {
	page := 0
	allOrdersFound := false
	ordersCount := 0

	for !allOrdersFound {
		// Retrieve orders from buyer account
		buyerOrders, err := getBuyerShippedOrders(page)
		if err != nil {
			logger.Error(fmt.Sprintf("failed to get orders on page :: %d", page), err)
			return
		}

		ordersCount = ordersCount + len(buyerOrders)
		if len(buyerOrders) == 0 {
			logger.Info(fmt.Sprintf("All orders have been found [%d]", ordersCount))
			allOrdersFound = true
			continue
		}

		for _, buyerOrder := range buyerOrders {
			// Fetch the order
			order, exists, err := getSellerOrderWithSellerOrderCode(buyerOrder.BuyerOrderCode)
			if err != nil {
				logger.Error("failed to get order with buyer order code", err)
				continue
			}

			if !exists {
				logger.Error("Order has not been synced to seller account", errors.New("error: order missing"))
				continue
			}

			// Check if the order has been marked as shipped on the seller account (retail side)
			//if hasBuyerOrderShipped(buyerOrder) && !order.Shipped {
			if buyerOrder.Shipped && !order.Shipped {
				logger.Info("Order has been shipped in buyer account, sharing it with the seller account")

				err := createFulfillmentOnSellerOrder(order.ID, buyerOrder.Fulfillments)
				if err != nil {
					logger.Error("failed to create fulfillment on the seller order", err)
					continue
				}
				logger.Info("Order has been marked as shipped in both accounts")
			} else if !buyerOrder.Shipped && order.Shipped {
				logger.Error("Order was marked as shipped in seller account but not buyer account", errors.New("invalid state"))
				continue
			} // Else: Shipped in both, or not shipped
		}

		page++
	}
}

// syncNewOrders :: Syncs any new orders from the seller account (retailer side) to the buyer account (supplier side)
func syncNewOrders() {
	page := 0
	allOrdersFound := false
	ordersCount := 0

	for !allOrdersFound {
		orders, err := getSellerNonShippedOrders(page)
		if err != nil {
			logger.Error(fmt.Sprintf("failed to get orders on page %d", page), err)
			return
		}
		ordersCount = ordersCount + len(orders)
		if len(orders) == 0 {
			logger.Info(fmt.Sprintf("All new orders have been found and synced [%d]", ordersCount))
			allOrdersFound = true
			continue
		}

		for _, order := range orders {
			// Check if exist on buyer/supplier side using the seller order code against the buyer order code
			_, exists, err := getBuyerOrderWithBuyerOrderCode(order.SellerOrderCode)
			if err != nil {
				logger.Error("failed to get order with buyer order code", err)
				continue
			}
			fmt.Printf("order :: %+v\n", order)

			if !exists {
				// Create new instance of the order on the buyer side
				buyerOrder, err := ConvertToBuyerOrder(order)
				if err != nil {
					logger.Error(fmt.Sprintf("Failed to convert order to buyer order for %s (Seller Order ID)", order.ID), err)
					continue
				}
				buyerOrderID, err := postNewBuyerOrderToAPI(buyerOrder)
				if err != nil {
					logger.Error(fmt.Sprintf("Failed to create new order for %s (Seller Order ID)", order.ID), err)
					continue
				}
				logger.Info(fmt.Sprintf("New order created on the buyer account :: %s --> %s", order.ID, buyerOrderID))
			}
		}

		page++
	}
}

func createFulfillmentOnSellerOrder(orderID string, fulfillments []Fulfillment) error {
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
		fmt.Printf("jsonPayload :: %+v\n", string(jsonPayload))

		_, err = http.PostRequest(fmt.Sprintf("/orders/%s/fulfillments", orderID), env.GetSellerAPIKey(), jsonPayload)
		if err != nil {
			return err
		}
		time.Sleep(time.Second * 1)
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

// getBuyerShippedOrders :: Returns a list of buyer orders that have been shipped
func getBuyerShippedOrders(page int) ([]Order, error) {
	resp, err := http.GetRequest("/orders?shipped=true", page, env.GetBuyerAPIKey())
	if err != nil {
		return []Order{}, err
	}

	var response []Order
	err = json.Unmarshal(resp, &response)
	if err != nil {
		return []Order{}, err
	}
	return response, nil
}

// getSellerOrderWithSellerOrderCode :: Returns a seller order using the seller order code from the seller API
func getSellerOrderWithSellerOrderCode(orderCode string) (Order, bool, error) {
	resp, err := http.GetRequest(fmt.Sprintf("/orders?sellerOrderCode=%s", orderCode), 0, env.GetSellerAPIKey())
	if err != nil {
		return Order{}, false, err
	}

	var response []Order
	err = json.Unmarshal(resp, &response)
	if err != nil {
		return Order{}, false, err
	}
	if len(response) == 0 {
		return Order{}, false, err
	} else if len(response) > 1 {
		return Order{}, false, errors.New(fmt.Sprintf("duplicate orders with order code (Under seller) :: %s", orderCode))
	}
	return response[0], true, nil
}

// getSellerNonShippedOrders :: Returns a list of (seller) orders that have not shipped from the seller API
func getSellerNonShippedOrders(page int) ([]Order, error) {
	resp, err := http.GetRequest("/orders?shipped=false", page, env.GetSellerAPIKey())
	if err != nil {
		return []Order{}, err
	}

	var response []Order
	err = json.Unmarshal(resp, &response)
	if err != nil {
		return []Order{}, err
	}
	return response, nil
}

// getBuyerOrderWithBuyerOrderCode :: Returns a buyer order from the buyer account using the list all orders endpoint
// and filter by the buyerOrderCode
// TODO - Using a seller get orders endpoint (should be buyer but it does not exist)
func getBuyerOrderWithBuyerOrderCode(buyerOrderCode string) (BuyerOrder, bool, error) {
	resp, err := http.GetRequest(fmt.Sprintf("/orders?buyerOrderCode=%s", buyerOrderCode), 0, env.GetBuyerAPIKey())
	if err != nil {
		return BuyerOrder{}, true, err
	}

	var response []Order
	err = json.Unmarshal(resp, &response)
	if err != nil {
		return BuyerOrder{}, true, err
	}
	if len(response) == 0 {
		return BuyerOrder{}, false, nil
	}
	return BuyerOrder{}, true, nil
}

// postNewBuyerOrderToAPI :: Submits a new order to the Buyer API for the buyer account
func postNewBuyerOrderToAPI(buyerOrder BuyerOrder) (string, error) {
	fmt.Printf("buyerOrder :: %+v\n", buyerOrder)
	jsonPayload, err := json.Marshal(buyerOrder)
	if err != nil {
		return "", err
	}

	resp, err := http.PostRequest("/buyer/orders", env.GetBuyerAPIKey(), jsonPayload)
	if err != nil {
		return "", err
	}

	var response BuyerOrder
	err = json.Unmarshal(resp, &response)
	if err != nil {
		return "", err
	}
	return response.ID, nil
}

// ConvertToBuyerOrder :: Converts an order from the seller order model to the buyer order model
func ConvertToBuyerOrder(o Order) (BuyerOrder, error) {
	buyerItems := []BuyerItem{}
	for _, item := range o.Items {
		// Look up the ID of the variant
		idOfVariant, err := products.GetIDOfVariantBySellerVariantCode(env.GetBuyerAPIKey(), item.SellerVariantCode)
		if err != nil {
			return BuyerOrder{}, err
		}
		buyerItems = append(buyerItems, BuyerItem{
			VariantID: idOfVariant,
			BuyerReference: item.ID,
			Quantity: item.Quantity,
		})
		time.Sleep(time.Second * 1)
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