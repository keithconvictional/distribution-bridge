package products

import (
	"distribution-bridge/alerts"
	"distribution-bridge/env"
	"distribution-bridge/global"
	"distribution-bridge/logger"
	"errors"
	"fmt"
	"time"
)

type Job struct {
	ID string
	Since *time.Time
	RequestManager *global.RequestManager
}


// SyncProducts :: Sync products from buyer account to seller account.
func (j *Job) SyncProducts() {
	page := 0
	allProductsFound := false
	productCount := 0
	// Fetch all products from seller accounts
	for !allProductsFound {
		products, err := j.getProductsFromAPI(page, env.GetBuyerAPIKey(), j.Since) // This doesn't actually filter as of 2021-09-17
		if err != nil {
			logger.Error(j.ID, global.DomainProducts, fmt.Sprintf("failed to get products on page %d", page), err)
			return
		}
		productCount = productCount + len(products)
		if len(products) == 0 {
			logger.Info(j.ID, global.DomainProducts, fmt.Sprintf("All products have been found [%d]", productCount))
			allProductsFound = true
			continue
		}
		logger.Info(j.ID, global.DomainProducts, fmt.Sprintf("%d products found (from buyer) on page %d", len(products), page))

		// For each product, it's consider to be new or exist on the buyer account
		for _, product := range products {
			j.SyncProduct(product)
			page++
		}
	}
}

// SyncProduct :: Sync a single product from buyer account to seller account
func (j *Job) SyncProduct(buyerProduct Product) {
	// Get the buyer product
	sellerProduct, exists, err := j.getProductFromAPIUsingCode(buyerProduct.Code, env.GetSellerAPIKey())
	if err != nil {
		logger.Error(j.ID, global.DomainProducts, fmt.Sprintf("failed to get product [%s]", buyerProduct.ID), err)
		return
	}

	//fmt.Printf("\n\n\n\n")  // TODO KEITH
	//fmt.Printf("Buyer Product :: %+v :: %+v\n\n", exists, sellerProduct)  // TODO KEITH
	//fmt.Printf("Seller Product :: %+v\n", buyerProduct)  // TODO KEITH
	//fmt.Printf("\n\n\n\n")  // TODO KEITH
	//
	//

	// If product exists
	validationErr := j.ValidateProduct(buyerProduct)
	if exists && validationErr != nil {
		// Delete the existing product from the buyer account since it's not valid
		logger.Error(j.ID, global.DomainProducts, fmt.Sprintf("Invalid seller product [https://app.convictional.com/products/%s] and buyer products [%s] because %s", buyerProduct.ID, sellerProduct.ID, validationErr.Error()), validationErr)
		alerts.SendAlert(fmt.Sprintf("Invalid seller product [https://app.convictional.com/products/%s] and buyer products [%s] because %s", buyerProduct.ID, sellerProduct.ID, validationErr.Error()))
		return
	} else if validationErr != nil {
		// Delete the existing product from the buyer account since it's not valid
		logger.Error(j.ID, global.DomainProducts, fmt.Sprintf("Invalid seller product [https://app.convictional.com/products/%s] because %s", buyerProduct.ID, validationErr.Error()), validationErr)
		alerts.SendAlert(fmt.Sprintf("Invalid seller product [https://app.convictional.com/products/%s] because %s", buyerProduct.ID, validationErr.Error()))
		return
	}

	// Apply PIM Updates
	// TODO
	logger.Info(j.ID, global.DomainProducts, "Apply to PIM")

	if exists {
		j.UpdateSellerProductFromBuyerProduct(sellerProduct, buyerProduct)
		return
	}
	j.CreateSellerProductFromBuyerProduct(buyerProduct)
}

// UpdateSellerProductFromBuyerProduct checks if the product has changed
func (j *Job) UpdateSellerProductFromBuyerProduct(sellerProduct Product, buyerProduct Product) {
	logger.Info(j.ID, global.DomainProducts, fmt.Sprintf("Product [%s] exists and checking for updates.", sellerProduct.Code))
	// Check changes, then update
	err := productsMatch(sellerProduct, buyerProduct)
	if err == nil {
		logger.Info(j.ID, global.DomainProducts, fmt.Sprintf("Products match between %s and %s", buyerProduct.ID, sellerProduct.ID))
		return
	}
	logger.Info(j.ID, global.DomainProducts, fmt.Sprintf("Products did not match between %s and %s b/c %+v", buyerProduct.ID, sellerProduct.ID, err))

	for _, image := range sellerProduct.Images {
		logger.Info(j.ID, global.DomainProducts, fmt.Sprintf("Delete image [%s] on product [%s]", image.ID, sellerProduct.ID))
		err = j.deleteProductImages(env.GetSellerAPIKey(), sellerProduct.ID, image.ID)
		if err != nil {
			logger.Error(j.ID, global.DomainProducts, fmt.Sprintf("Failed to delete product images [%s]", sellerProduct.ID), err)
			return
		}
	}

	// Temp save seller product ID
	body := ProductUpdateBody{
		Title: &buyerProduct.Title,
		Active: &buyerProduct.Active,
		BodyHTML: &buyerProduct.BodyHTML,
		Images: RemoveImagesID(buyerProduct.Images),
		Tags: &buyerProduct.Tags,
		Vendor: &buyerProduct.Vendor,
		Variants: RemoveVariantsIDs(buyerProduct.Variants),
		Options: RemoveOptionsIDs(buyerProduct.Options),
		Attributes: &buyerProduct.Attributes,
	}

	fmt.Printf("body :: %+v\n", body)  // TODO KEITH
	fmt.Printf("sellerProduct.ID :: %+v\n", sellerProduct.ID)  // TODO KEITH


	// Mark updated product as inactive
	if env.ProductUpdatesToInActive() {
		active := false
		body.Active = &active
	}
	err = j.updateProductOnAPI(env.GetSellerAPIKey(), sellerProduct.ID, body)
	if err != nil {
		logger.Error(j.ID, global.DomainProducts, fmt.Sprintf("failed to update the product on seller account (Existing) :: %s", sellerProduct.ID), err)
		return
	}

	// Update images :: Image src fields cannot be updated base on API docs, so all images must be delete then re-added.
}

func RemoveImagesID(images []Image) *[]Image {
	for i, _ := range images {
		images[i].ID = ""
	}
	return &images
}

func RemoveVariantsIDs(variants []Variant) *[]Variant {
	for i, _ := range variants {
		variants[i].ID = ""
	}
	return &variants
}

func RemoveOptionsIDs(options []Option) *[]Option {
	for i, _ := range options {
		options[i].ID = ""
	}
	return &options
}

// CreateSellerProductFromBuyerProduct creates a new seller product from the buyer product
func (j *Job) CreateSellerProductFromBuyerProduct(sellerProduct Product) {
	if err := j.ValidateProduct(sellerProduct); err != nil {
		logger.Error(j.ID, global.DomainProducts, fmt.Sprintf("Product failed to create because [%s] because [%s]", sellerProduct.ID, err.Error()), err)
		return
	}
	logger.Info(j.ID, global.DomainProducts, fmt.Sprintf("Product [%s] does not exist and creating new instance.", sellerProduct.Code))
	// Create new product on buyer account
	productID, err := j.createProductOnAPI(sellerProduct, env.GetSellerAPIKey())
	if err != nil {
		logger.Error(j.ID, global.DomainProducts, fmt.Sprintf("failed to create new product on seller account :: Seller Product ID [%s]", sellerProduct.ID), err)
		// Not supported but push error to the seller product
		return
	}

	logger.Info(j.ID, global.DomainProducts, fmt.Sprintf("New product created on seller account :: %s --> %s", sellerProduct.ID, productID))

	// Mark new product as inactive
	if env.NewProductToInActive() {
		err = j.updateProductAsInactiveOnAPI(env.GetSellerAPIKey(), productID)
		if err != nil {
			logger.Error(j.ID, global.DomainProducts,fmt.Sprintf("failed to mark product as inactive on seller account (New) :: %s", productID), err)
		}
	}
}

// ValidateProduct required data
func (j *Job) ValidateProduct(product Product) error {
	if product.Title == "" {
		return errors.New("error: no title")
	}
	for _, variant := range product.Variants {
		if variant.Barcode == "" {
			return fmt.Errorf("error: no barcode on %s [%s]", product.Title, product.ID)
		}
		if variant.BarcodeType == "" {
			return fmt.Errorf("no barcode type on %s [%s]", product.Title, product.ID)
		}
	}
	return nil
}