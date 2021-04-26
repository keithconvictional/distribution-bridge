package orders

import (
	"distribution-bridge/env"
	"distribution-bridge/http"
	"distribution-bridge/logger"
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
			// Check if exist on buyer/supplier side
			//order, exists, err := getOrderWithBuyerOrderCode
			fmt.Printf("order :: %+v\n", order)

			// Create new instance of the order on the buyer side
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