package orders

import (
	"distribution-bridge/env"
	"distribution-bridge/http"
	"distribution-bridge/logger"
	"distribution-bridge/products"
	"encoding/json"
	"fmt"
)

// Sync new orders from buyer account to seller account. Sync order updates both ways.
func SyncOrders() {
	// Get new orders from seller account (Retailer side)
	syncNewOrders()

	// Get order updates from buyer account (Supplier side)
}

func syncNewOrders() {
	page := 0
	allOrdersFound := false
	orderCount := 0

	for !allOrdersFound {
		orders, err := getOrdersFromAPI(page, env.GetSellerAPIKey())
		if err != nil {
			logger.Error(fmt.Sprintf("failed to get orders on page %d", page), err)
			return
		}
		orderCount = orderCount + len(orders)
		if len(orders) == 0 {
			logger.Info(fmt.Sprintf("All orders have been found [%d]", orderCount))
			allOrdersFound = true
			continue
		}

		for _, order := range orders {
			// Check if exist on buyer/supplier side using the seller order code against the buyer order code
			_, exists, err := getBuyerOrderWithBuyerOrderCode(env.GetBuyerAPIKey(), order.SellerOrderCode)
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
				buyerOrderID, err := postNewBuyerOrderToAPI(env.GetBuyerAPIKey(), buyerOrder)
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

func getOrdersFromAPI(page int, apiKey string) ([]Order, error) {
	resp, err := http.GetRequest("/orders", page, apiKey)
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

// Using a seller get orders endpoint (should be buyer but it does not exist)
func getBuyerOrderWithBuyerOrderCode(apiKey string, buyerOrderCode string) (BuyerOrder, bool, error) {
	resp, err := http.GetRequest(fmt.Sprintf("/orders?buyerOrderCode=%s", buyerOrderCode), 0, apiKey)
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

func postNewBuyerOrderToAPI(apiKey string, buyerOrder BuyerOrder) (string, error) {
	fmt.Printf("buyerOrder :: %+v\n", buyerOrder)
	jsonPayload, err := json.Marshal(buyerOrder)
	if err != nil {
		return "", err
	}

	resp, err := http.PostRequest("/buyer/orders", apiKey, jsonPayload)
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