package main

import (
	"distribution-bridge/env"
	"distribution-bridge/logger"
	"distribution-bridge/orders"
	"fmt"
)

func main() {
	fmt.Println("Starting Distribution Bridge...")
	// Check for variables
	if !env.ValidEnvVariables() {
		logger.Info("Required environment variables are missing")
		return
	}

	// Sync products
	//products.SyncProducts()

	// Sync orders
	if env.DropShippingEnabled() {
		logger.Info("Drop shipping is enabled.")
		orders.SyncOrders()
	}
}



