package products

import "time"

// Created with: https://mholt.github.io/json-to-go/
type Product struct {
	ID              string     `json:"_id"`
	Code            string     `json:"code"`
	Active          bool       `json:"active"`
	BodyHTML        string     `json:"bodyHtml"`
	Images          []Image   `json:"images"`
	Tags            []string   `json:"tags"`
	Title           string     `json:"title"`
	Vendor          string     `json:"vendor"`
	Variants        []Variant `json:"variants"`
	Options         []Option  `json:"options"`
	DelistedUpdated time.Time  `json:"delistedUpdated"`
	Created         time.Time  `json:"created"`
	Updated         time.Time  `json:"updated"`
	CompanyObjectID string     `json:"companyObjectId"`
	Type            string     `json:"type"`
	CompanyID       string     `json:"companyId"`
	Attributes map[string]string `json:"attributes"`
}

type Image struct {
	ID         string        `json:"_id,omitempty"`
	Src        string        `json:"src"`
	Position   int           `json:"position"`
	VariantIds []interface{} `json:"variantIds"`
}
type Dimensions struct {
	Length int    `json:"length"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Units  string `json:"units"`
}
type Variant struct {
	ID                string     `json:"_id,omitempty"`
	Title             string     `json:"title"`
	RetailPrice       float64    `json:"retailPrice"`
	InventoryQuantity int        `json:"inventory_quantity"`
	SkipCount         bool       `json:"skipCount"`
	Weight            int        `json:"weight"`
	WeightUnits       string     `json:"weightUnits"`
	Dimensions        Dimensions `json:"dimensions"`
	Sku               string     `json:"sku"`
	Barcode           string     `json:"barcode"`
	BarcodeType       string     `json:"barcodeType"`
	Code              string     `json:"code"`
	VariantID         int        `json:"id"`
	Option1           string     `json:"option1"`
	Option2           string     `json:"option2"`
	Option3           string     `json:"option3"`
}
type Option struct {
	ID       string `json:"_id,omitempty"`
	Name     string `json:"name"`
	Position int    `json:"position"`
	Type     string `json:"type"`
}


type ProductUpdateBody struct {
	Title *string `json:"title"`
	Active *bool `json:"active"`
	BodyHTML *string `json:"bodyHtml"`
	Tags *[]string `json:"tags"`
	Vendor *string `json:"vendor"`
	Variants *[]Variant `json:"variants"`
	Images *[]Image `json:"image"`
	Options         *[]Option  `json:"options"`
	Attributes *map[string]string `json:"attributes"`
}