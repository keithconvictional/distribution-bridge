package products

import (
	"distribution-bridge/global"
	"distribution-bridge/http"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// getProductsFromAPI calls the get products endpoint
func (j *Job) getProductsFromAPI(page int, apiKey string, since *time.Time) ([]Product, error) {
	j.RequestManager.Wait()

	resp, err := http.GetRequest(j.ID, global.DomainProducts,"/products", page, apiKey, j.Since)
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


func (j *Job) getProductFromAPIUsingCode(code string, apiKey string) (Product, bool, error) {
	j.RequestManager.Wait()

	resp, err := http.GetRequest(j.ID, global.DomainProducts, fmt.Sprintf("/products?productCode=%s", code), 0, apiKey, j.Since)
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


func (j *Job) createProductOnAPI(product Product, apiKey string) (string, error) {
	j.RequestManager.Wait()

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

func (j *Job) updateProductOnAPI(apiKey string, productID string, body ProductUpdateBody) error {
	j.RequestManager.Wait()

	jsonPayload, err := json.Marshal(body)
	if err != nil {
		return err
	}

	fmt.Printf("jsonPayload :: %s\n", string(jsonPayload)) // TODO KEITH

	resp, err := http.PatchRequest(fmt.Sprintf("/products/%s", productID), apiKey, jsonPayload)
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

func (j *Job) deleteProductImages(apiKey string, productID string, imageID string) error {
	j.RequestManager.Wait()

	resp, err := http.DeleteRequest(fmt.Sprintf("/products/%s/images/%s", productID, imageID), apiKey)
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

func (j *Job) updateProductAsInactiveOnAPI(apiKey string, productID string) error {
	active := false
	return j.updateProductOnAPI(apiKey, productID, ProductUpdateBody{
		Active: &active,
	})
}

// Buyer account calling a seller endpoint (This should be fixed)
func (j *Job) GetIDOfVariantBySellerBarcode(apiKey string, barcode string, barcodeType string) (string, error) {
	page := 0
	found := false
	for !found {
		j.RequestManager.Wait()

		resp, err := http.GetRequest(j.ID, global.DomainProducts, "/products", page, apiKey, nil)
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

		fmt.Printf("products :: %+v\n", products)
		for _, product := range products {
			for _, variant := range product.Variants {
				if variant.Barcode == barcode && variant.BarcodeType == barcodeType {
					return variant.ID, nil
				}
			}
		}

		page++
	}
	return "", errors.New(fmt.Sprintf("ID of variant not found using barcode (%s / %s)", barcodeType, barcode))
}