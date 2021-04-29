package products

import (
	"distribution-bridge/env"
	"distribution-bridge/http"
	"distribution-bridge/logger"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"time"
)



// Sync products from seller account to buyer account.
func SyncProducts() {
	page := 0
	allProductsFound := false
	productCount := 0
	// Fetch all products from seller accounts
	for !allProductsFound {
		products, err := getProductsFromAPI(page, env.GetBuyerAPIKey())
		if err != nil {
			logger.Error(fmt.Sprintf("failed to get products on page %d", page), err)
			return
		}
		productCount = productCount + len(products)
		if len(products) == 0 {
			logger.Info(fmt.Sprintf("All products have been found [%d]", productCount))
			allProductsFound = true
			continue
		}

		// For each product, it's consider to be new or exist on the buyer account
		for _, product := range products {
			sellerProduct, exists, err := getProductFromAPIUsingCode(product.Code, env.GetSellerAPIKey())
			if err != nil {
				logger.Error(fmt.Sprintf("failed to get product [%s]", product.ID), err)
				return
			}

			// Apply PIM updates (This would be any configured overwrites that have been setup)
			// Not built :: Ex. Loops through Google sheet and when product code = product.Code then updates with columns for corresponding data. Or PIM provider, we make an outbound call to them and they return the updated product

			if exists {
				logger.Info(fmt.Sprintf("Product [%s] exists and checking for updates.", product.Code))
				// Check changes, then update
				err := productsMatch(product, sellerProduct)
				if err == nil {
					logger.Info(fmt.Sprintf("Products match between %s and %s", product.ID, sellerProduct.ID))
					continue
				} else {
					logger.Info(fmt.Sprintf("Products did not match between %s and %s b/c %+v", product.ID, sellerProduct.ID, err))
					// Temp save seller product ID
					sellerProductID := sellerProduct.ID
					sellerProduct = product
					sellerProduct.ID = sellerProductID
					// Mark updated product as inactive
					if env.ProductUpdatesToInActive() {
						sellerProduct.Active = false
					}
					err = updateProductOnAPI(env.GetSellerAPIKey(), sellerProduct)
					if err != nil {
						logger.Error(fmt.Sprintf("failed to update the product on seller account (Existing) :: %s", sellerProduct.ID), err)
					}
				}
			} else {
				logger.Info(fmt.Sprintf("Product [%s] does not exist and creating new instance.", product.Code))
				// Create new product on buyer account
				productID, err := createProductOnAPI(product, env.GetSellerAPIKey())
				if err != nil {
					logger.Error(fmt.Sprintf("failed to create new product on seller account :: Seller Product ID [%s]", product.ID), err)
					// Not supported but push error to the seller product
					return
				}
				logger.Info(fmt.Sprintf("New product created on buyer account :: %s --> %s", product.ID, productID))

				// Mark new product as inactive
				if env.NewProductToInActive() {
					product.Active = false
					sellerProduct, _, err = getProductFromAPIUsingCode(product.Code, env.GetSellerAPIKey())
					if err != nil {
						logger.Error(fmt.Sprintf("failed to get product for seller [%s]", product.ID), err)
						return
					}

					err := updateProductOnAPI(env.GetSellerAPIKey(), sellerProduct)
					logger.Error(fmt.Sprintf("failed to mark product as inactive on seller account (New) :: %s", sellerProduct.ID), err)
					// Not supported but push error to the buyer and seller product
					return
				}
			}

			time.Sleep(time.Second * 1) // Max of 4 API calls in this block, so this is bad rate limiting
		}

		page++
	}
}

// productsMatch custom method for comparing two products. IDs will be completely different in both.
func productsMatch(product Product, productTwo Product) error {
	// Images
	if len(product.Images) != len(productTwo.Images) {
		return errors.New("unequal number of images between both products")
	}
	foundSrcs := 0
	for _, imageFromOne := range product.Images {
		for _, imageFromTwo := range productTwo.Images {
			if imageFromOne.Src == imageFromTwo.Src {
				foundSrcs++

				if imageFromOne.Position != imageFromTwo.Position {
					return errors.New(fmt.Sprintf("image positions do not match for %s", imageFromOne.Src))
				}
				break
			}
		}
	}
	if foundSrcs != len(product.Images) {
		return errors.New(fmt.Sprintf("did not find all matches for all images (Found %d of %d)", foundSrcs, len(product.Images)))
	}

	// Variants
	if len(product.Variants) != len(productTwo.Variants) {
		return errors.New("unequal number of variants")
	}
	foundVariants := 0
	for _, variantOne := range product.Variants {
		for _, variantTwo := range product.Variants {
			if variantOne.VariantID == variantTwo.VariantID {
				foundVariants++
				if !cmp.Equal(variantOne, variantTwo, cmpopts.IgnoreFields(Variants{}, "ID")) {
					return errors.New(fmt.Sprintf("variants do not match for %s and %s", variantOne.ID, variantTwo.ID))
				}
				break
			}
		}
	}
	if foundVariants != len(product.Variants) {
		return errors.New(fmt.Sprintf("did not find all variants (Found %d of %d)", foundVariants, len(product.Variants)))
	}

	// Options
	if len(product.Options) != len(productTwo.Options) {
		return errors.New("unequal number of options")
	}
	foundOptions := 0
	for _, optionOne := range product.Options {
		for _, optionTwo := range product.Options {
			if optionOne.Name == optionTwo.Name {
				foundOptions++
				if !cmp.Equal(optionOne, optionTwo, cmpopts.IgnoreFields(Options{}, "ID")) {
					return errors.New(fmt.Sprintf("did not find matching options between %s and %s", optionOne.Name, optionTwo.Name))
				}
				break
			}
		}
	}
	if foundOptions != len(product.Options) {
		return errors.New(fmt.Sprintf("did not find all options (Found %d of %d)", foundOptions, len(product.Options)))
	}

	// All other fields that should match
	if !cmp.Equal(product, productTwo, cmpopts.IgnoreFields(Product{}, "ID", "Active", "Images", "Variants", "Options", "Created", "Updated", "CompanyObjectID","CompanyID")) {
		return errors.New("products do not match")
	}

	return nil
}

// getProductsFromAPI calls the get products endpoint
func getProductsFromAPI(page int, apiKey string) ([]Product, error) {
	resp, err := http.GetRequest("/products", page, apiKey)
	if err != nil {
		return []Product{}, err
	}

	var response []Product
	err = json.Unmarshal(resp, &response)
	if err != nil {
		return []Product{}, err
	}
	return response, nil
}

func getProductFromAPIUsingCode(code string, apiKey string) (Product, bool, error) {
	resp, err := http.GetRequest(fmt.Sprintf("/products?productCode=%s", code), 0, apiKey)
	if err != nil {
		return Product{}, false, err
	}

	// Capture 404
	fmt.Printf("getProductFromAPI :: err :: %+v\n", err)

	var response []Product
	err = json.Unmarshal(resp, &response)
	if err != nil {
		return Product{}, false, err
	}
	if len(response) == 0 {
		return Product{}, false, nil
	}
	return response[0], true, nil
}

func createProductOnAPI(product Product, apiKey string) (string, error) {
	fmt.Printf("product :: %+v", product)
	jsonPayload, err := json.Marshal(product)
	if err != nil {
		return "", err
	}

	resp, err := http.PostRequest("/products", apiKey, jsonPayload)
	if err != nil {
		return "", err
	}

	var response Product
	err = json.Unmarshal(resp, &response)
	if err != nil {
		return "", err
	}
	return product.ID, nil
}

func updateProductOnAPI(apiKey string, product Product) error {
	jsonPayload, err := json.Marshal(product)
	if err != nil {
		return err
	}

	resp, err := http.PutRequest(fmt.Sprintf("/products/%s", product.ID), apiKey, jsonPayload)
	if err != nil {
		return err
	}

	var response Product
	err = json.Unmarshal(resp, &response)
	if err != nil {
		return err
	}
	return nil
}

// Buyer account calling a seller endpoint (This should be fixed)
func GetIDOfVariantBySellerVariantCode(apiKey string, sellerVariantCode string) (string, error) {
	page := 0
	found := false
	for !found {
		resp, err := http.GetRequest("/products", page, apiKey)
		if err != nil {
			return "", err
		}
		var products []Product
		err = json.Unmarshal(resp, &products)
		if err != nil {
			return "", err
		}

		if len(products) == 0 {
			found = true
		}

		for _, product := range products {
			for _, variant := range product.Variants {
				if variant.Code == sellerVariantCode {
					return variant.ID, nil
				}
			}
		}

		page++
		time.Sleep(time.Second * 1)
	}
	return "", errors.New(fmt.Sprintf("ID of variant not found using variantID/Code (%s)", sellerVariantCode))
}