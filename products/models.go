package products

import "time"

// Created with: https://mholt.github.io/json-to-go/
type Product struct {
	ID              string     `json:"_id"`
	Code            string     `json:"code"`
	Active          bool       `json:"active"`
	BodyHTML        string     `json:"bodyHtml"`
	Images          []Images   `json:"images"`
	Tags            []string   `json:"tags"`
	Title           string     `json:"title"`
	Vendor          string     `json:"vendor"`
	Variants        []Variants `json:"variants"`
	Options         []Options  `json:"options"`
	DelistedUpdated time.Time  `json:"delistedUpdated"`
	Created         time.Time  `json:"created"`
	Updated         time.Time  `json:"updated"`
	CompanyObjectID string     `json:"companyObjectId"`
	Type            string     `json:"type"`
	CompanyID       string     `json:"companyId"`
}

type Images struct {
	ID         string        `json:"_id"`
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
type Variants struct {
	ID                string     `json:"_id"`
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
type Options struct {
	ID       string `json:"_id"`
	Name     string `json:"name"`
	Position int    `json:"position"`
	Type     string `json:"type"`
}