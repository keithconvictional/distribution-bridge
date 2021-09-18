package main

import (
	"distribution-bridge/env"
	"distribution-bridge/global"
	"distribution-bridge/logger"
	"distribution-bridge/orders"
	"distribution-bridge/products"
	"github.com/google/uuid"
	"time"
)

func main() {
	jobID := uuid.New().String()
	logger.Info(jobID, global.DomainGeneral, "Starting Distribution Bridge...")
	// Check for variables
	if !env.ValidEnvVariables(jobID) {
		logger.Info(jobID, global.DomainGeneral, "Required environment variables are missing")
		return
	}

	var since *time.Time
	requestManager := global.RequestManager{}

	// Sync products
	productJob := products.Job{
		ID: jobID,
		Since: since,
		RequestManager: &requestManager,
	}
	productJob.SyncProducts()

	// Sync orders
	ordersJob := orders.Job{
		ID: jobID,
		Since: since,
		ProductsJob: productJob,
		RequestManager: &requestManager,
	}
	ordersJob.SyncOrders()
}



