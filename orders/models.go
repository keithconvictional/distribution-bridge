package orders

import (
	"time"
)

type Order struct {
	ID               string `json:"_id"`
	BuyerOrderCode   string `json:"buyerOrderCode"`
	SellerOrderCode  string `json:"sellerOrderCode"`
	BuyerCompanyID   string `json:"buyerCompanyId"`
	SellerCompanyID  string `json:"sellerCompanyId"`
	BuyerEmail       string `json:"buyerEmail"`
	Currency         string `json:"currency"`
	InvoiceID        string `json:"invoiceId"`
	Note             string `json:"note"`
	HasCancellations bool   `json:"hasCancellations"`
	ShippingAddress  struct {
		Name       string `json:"name"`
		AddressOne string `json:"addressOne"`
		AddressTwo string `json:"addressTwo"`
		City       string `json:"city"`
		State      string `json:"state"`
		Country    string `json:"country"`
		Zip        string `json:"zip"`
		Company    string `json:"company"`
	} `json:"shippingAddress"`
	FillTime float64 `json:"fillTime"`
	ShipTime float64 `json:"shipTime"`
	Custom   []struct {
		Key   string `json:"key"`
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"custom"`
	Items []struct {
		ID                string    `json:"_id"`
		SellerVariantCode string    `json:"sellerVariantCode"`
		Quantity          int       `json:"quantity"`
		Cancelled         bool      `json:"cancelled"`
		CancelledReason   string    `json:"cancelledReason"`
		CancelledDate     time.Time `json:"cancelledDate"`
	} `json:"items"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
	Posted      bool      `json:"posted"`
	PostedDate  time.Time `json:"postedDate"`
	Shipped     bool      `json:"shipped"`
	ShippedDate time.Time `json:"shippedDate"`
	Billed      bool      `json:"billed"`
	BilledDate  time.Time `json:"billedDate"`
	Fulfillments []Fulfillment `json:"fulfillments"`
}

type Fulfillment struct {
	ID                    string   `json:"_id"`
	BuyerFulfillmentCode  string   `json:"buyerFulfillmentCode"`
	SellerFulfillmentCode string   `json:"sellerFulfillmentCode"`
	Carrier               string   `json:"carrier"`
	TrackingCode          string   `json:"trackingCode"`
	TrackingUrls          []string `json:"trackingUrls"`
	Items                 []struct {
		Quantity          int           `json:"quantity"`
		ID                string        `json:"_id"`
		OrderItemID       string        `json:"orderItemId"`
		BuyerItemCode     string        `json:"buyerItemCode"`
		SellerItemCode    string        `json:"sellerItemCode"`
		BuyerVariantCode  string        `json:"buyerVariantCode"`
		SellerVariantCode string        `json:"sellerVariantCode"`
		BuyerProductCode  string        `json:"buyerProductCode"`
		SellerProductCode string        `json:"sellerProductCode"`
		Type              string        `json:"type"`
		Title             string        `json:"title"`
		Sku               string        `json:"sku"`
		Price             int           `json:"price"`
		RetailPrice       int           `json:"retailPrice"`
		Barcode           string        `json:"barcode"`
		BarcodeType       string        `json:"barcodeType"`
		Weight            int           `json:"weight"`
		Custom            []interface{} `json:"custom"`
	} `json:"items"`
}

type NewFulfillmentRequestBody struct {
	Carrier string `json:"carrier"`
	TrackingCode string `json:"trackingCode"`
	TrackingURLs []string `json:"trackingUrls"`
	Items []NewFulfillmentItem `json:"items"`
}

type NewFulfillmentItem struct {
	ID int32 `json:"id"`
	SKU string `json:"sku"`
	Quantity int32 `json:"quantity"`
}


type BuyerOrder struct {
	ID             string    `json:"id"`
	BuyerReference string    `json:"buyerReference"`
	OrderedDate    time.Time `json:"orderedDate"`
	Created        time.Time `json:"created"`
	Updated        time.Time `json:"updated"`
	Address        struct {
		Name       string `json:"name"`
		AddressOne string `json:"addressOne"`
		AddressTwo string `json:"addressTwo"`
		City       string `json:"city"`
		State      string `json:"state"`
		Country    string `json:"country"`
		Zip        string `json:"zip"`
		Company    string `json:"company"`
	} `json:"address"`
	Items []BuyerItem `json:"items"`
	Note         string `json:"note"`
	SellerOrders []struct {
		ID              string    `json:"id"`
		BuyerOrderID    string    `json:"buyerOrderId"`
		BuyerReference  string    `json:"buyerReference"`
		SellerReference string    `json:"sellerReference"`
		CompanyID       string    `json:"companyId"`
		BaseCurrency    string    `json:"baseCurrency"`
		PackingSlipURL  string    `json:"packingSlipUrl"`
		InvoiceID       string    `json:"invoiceId"`
		Posted          bool      `json:"posted"`
		PostedDate      string    `json:"postedDate"`
		Fulfilled       bool      `json:"fulfilled"`
		FulfilledDate   string    `json:"fulfilledDate"`
		Invoiced        bool      `json:"invoiced"`
		InvoicedDate    string    `json:"invoicedDate"`
		Refunded        bool      `json:"refunded"`
		RefundedDate    string    `json:"refundedDate"`
		Created         time.Time `json:"created"`
		Updated         time.Time `json:"updated"`
		Address         struct {
			Name       string `json:"name"`
			AddressOne string `json:"addressOne"`
			AddressTwo string `json:"addressTwo"`
			City       string `json:"city"`
			State      string `json:"state"`
			Country    string `json:"country"`
			Zip        string `json:"zip"`
			Company    string `json:"company"`
		} `json:"address"`
		Items []struct {
			ID        string  `json:"id"`
			VariantID string  `json:"variantId"`
			Quantity  int     `json:"quantity"`
			BasePrice float64 `json:"basePrice"`
		} `json:"items"`
		Fulfillments []struct {
			ID           string    `json:"id"`
			Posted       bool      `json:"posted"`
			PostedDate   string    `json:"postedDate"`
			Created      time.Time `json:"created"`
			Updated      time.Time `json:"updated"`
			Carrier      string    `json:"carrier"`
			TrackingCode string    `json:"trackingCode"`
			TrackingUrls []string  `json:"trackingUrls"`
			Items        []struct {
				ID          string `json:"id"`
				OrderItemID string `json:"orderItemId"`
				Quantity    int    `json:"quantity"`
			} `json:"items"`
		} `json:"fulfillments"`
	} `json:"sellerOrders"`
	Metafields *struct {
		AdditionalProp struct {
			AdditionalProp struct {
				Value       string    `json:"value"`
				Description string    `json:"description"`
				Updated     time.Time `json:"updated"`
			} `json:"additionalProp"`
		} `json:"additionalProp"`
	} `json:"metafields"`
}

type BuyerItem struct {
	ID                string  `json:"id"`
	VariantID         string  `json:"variantId"`
	BuyerReference    string  `json:"buyerReference"`
	SellerOrderID     string  `json:"sellerOrderId"`
	SellerOrderItemID string  `json:"sellerOrderItemId"`
	Quantity          int     `json:"quantity"`
	RetailPrice       float64 `json:"retailPrice"`
}