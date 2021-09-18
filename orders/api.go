package orders

import (
	"distribution-bridge/env"
	"distribution-bridge/global"
	"distribution-bridge/http"
	"encoding/json"
	"errors"
	"fmt"
)

// getBuyerOrderWithBuyerOrderCode :: Returns a buyer order from the buyer account using the list all orders endpoint
// and filter by the buyerOrderCode
// TODO - Using a seller get orders endpoint (should be buyer but it does not exist)
func (j *Job) getBuyerOrderWithBuyerOrderCode(buyerOrderCode string) (BuyerOrder, bool, error) {
	j.RequestManager.Wait()

	resp, err := http.GetRequest(j.ID, global.DomainOrders, fmt.Sprintf("/orders?buyerOrderCode=%s", buyerOrderCode), 0, env.GetBuyerAPIKey(), j.Since)
	if err != nil {
		return BuyerOrder{}, true, err
	}

	var response []Order
	err = json.Unmarshal(resp, &response)
	if err != nil {
		return BuyerOrder{}, true, err
	}
	if len(response) == 0 {
		return BuyerOrder{}, false, nil
	}
	return BuyerOrder{}, true, nil
}

// postNewBuyerOrderToAPI :: Submits a new order to the Buyer API for the buyer account
func (j *Job) postNewBuyerOrderToAPI(buyerOrder BuyerOrder) (string, error) {
	j.RequestManager.Wait()

	fmt.Printf("buyerOrder :: %+v\n", buyerOrder) // TODO KEITH
	jsonPayload, err := json.Marshal(buyerOrder)
	if err != nil {
		return "", err
	}

	resp, err := http.PostRequest("/buyer/orders", env.GetBuyerAPIKey(), jsonPayload)
	if err != nil {
		return "", err
	}

	var response BuyerOrder
	err = json.Unmarshal(resp, &response)
	if err != nil {
		return "", err
	}
	return response.ID, nil
}


// getBuyerShippedOrders :: Returns a list of buyer orders that have been shipped
func (j *Job) getBuyerShippedOrders(page int) ([]Order, error) {
	j.RequestManager.Wait()

	resp, err := http.GetRequest(j.ID, global.DomainOrders, "/orders?shipped=true", page, env.GetBuyerAPIKey(), nil)
	if err != nil {
		return []Order{}, err
	}

	var response []Order
	err = json.Unmarshal(resp, &response)
	if err != nil {
		return []Order{}, err
	}
	return response, nil
}

// getSellerOrderWithSellerOrderCode :: Returns a seller order using the seller order code from the seller API
func (j *Job) getSellerOrderWithSellerOrderCode(orderCode string) (Order, bool, error) {
	j.RequestManager.Wait()

	resp, err := http.GetRequest(j.ID, global.DomainOrders, fmt.Sprintf("/orders?sellerOrderCode=%s", orderCode), 0, env.GetSellerAPIKey(), nil)
	if err != nil {
		return Order{}, false, err
	}

	var response []Order
	err = json.Unmarshal(resp, &response)
	if err != nil {
		return Order{}, false, err
	}
	if len(response) == 0 {
		return Order{}, false, err
	} else if len(response) > 1 {
		return Order{}, false, errors.New(fmt.Sprintf("duplicate orders with order code (Under seller) :: %s", orderCode))
	}
	return response[0], true, nil
}

// getSellerNonShippedOrders :: Returns a list of (seller) orders that have not shipped from the seller API
func (j *Job) getSellerNonShippedOrders(page int) ([]Order, error) {
	j.RequestManager.Wait()

	resp, err := http.GetRequest(j.ID, global.DomainOrders, "/orders?shipped=false", page, env.GetSellerAPIKey(), nil)
	if err != nil {
		return []Order{}, err
	}

	var response []Order
	err = json.Unmarshal(resp, &response)
	if err != nil {
		return []Order{}, err
	}
	return response, nil
}
