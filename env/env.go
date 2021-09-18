package env

import (
	"distribution-bridge/global"
	"distribution-bridge/logger"
	"os"
	"strings"
)

func ValidEnvVariables(jobID string) bool {
	if GetSellerAPIKey() == "" {
		logger.Info(jobID, global.DomainGeneral, "No SELLER_API_KEY set")
		return false
	}
	if GetBuyerAPIKey() == "" {
		logger.Info(jobID, global.DomainGeneral, "No BUYER_API_KEY set")
		return false
	}
	return true
}

func GetSellerAPIKey() string {
	return os.Getenv("SELLER_API_KEY")
}

func GetBuyerAPIKey() string {
	return os.Getenv("BUYER_API_KEY")
}

func DropShippingEnabled() bool {
	return getEnvBool("DROP_SHIPPING_ENABLED", true)
}

func ProductUpdatesToInActive() bool {
	return getEnvBool("PRODUCT_UPDATES_TO_INACTIVE", false)
}

func NewProductToInActive() bool {
	return getEnvBool("NEW_PRODUCT_TO_INACTIVE", true)
}

func getEnvBool(key string, def bool) bool {
	str := os.Getenv(key)
	if str == "" {
		// Default
		return def
	}
	if strings.ToLower(str) == "false" {
		return false
	}
	return true
}

func GetBaseURL() string {
	baseURL := os.Getenv("CONVICTIONAL_API_URL")
	if baseURL != "" {
		return baseURL
	}
	return "https://api.convictional.com"
}