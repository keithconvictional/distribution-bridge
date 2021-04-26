package orders

import "time"

type Order struct {
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
	Items []struct {
		ID                string  `json:"id"`
		VariantID         string  `json:"variantId"`
		BuyerReference    string  `json:"buyerReference"`
		SellerOrderID     string  `json:"sellerOrderId"`
		SellerOrderItemID string  `json:"sellerOrderItemId"`
		Quantity          int     `json:"quantity"`
		RetailPrice       float64 `json:"retailPrice"`
	} `json:"items"`
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
	Metafields struct {
		AdditionalProp struct {
			AdditionalProp struct {
				Value       string    `json:"value"`
				Description string    `json:"description"`
				Updated     time.Time `json:"updated"`
			} `json:"additionalProp"`
		} `json:"additionalProp"`
	} `json:"metafields"`
}