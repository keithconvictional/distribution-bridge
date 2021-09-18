package http

import (
	"bytes"
	"distribution-bridge/env"
	"distribution-bridge/logger"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var httpClient = &http.Client{
	Timeout: time.Second * 10,
}

func GetRequest(requestID string, domain string, urlPath string, page int, apiKey string, since *time.Time) ([]byte, error) {
	if !strings.Contains(urlPath, "?") {
		urlPath += "?"
	} else {
		urlPath += "&"
	}
	url := fmt.Sprintf("%s%spage=%d&limit=250", env.GetBaseURL(), urlPath, page)
	if since != nil {
		url += fmt.Sprintf("&updated=%s", since.Format("2006-01-02T15:04:05Z"))
	}
	logger.Info(requestID, domain, fmt.Sprintf("Calling url :: %s\n", url))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []byte{}, err
	}

	// Add headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", apiKey)

	resp, err := sendRequest(req)
	if err != nil {
		return []byte{}, err
	}
	return resp, nil
}

func PostRequest(urlPath string, apiKey string, jsonPayload []byte) ([]byte, error) {
	return requestWithBody(urlPath, "POST", apiKey, jsonPayload)
}

func PatchRequest(urlPath string, apiKey string, jsonPayload []byte) ([]byte, error) {
	return requestWithBody(urlPath, "PATCH", apiKey, jsonPayload)
}

func DeleteRequest(urlPath string, apiKey string) ([]byte, error) {
	url := fmt.Sprintf("%s%s", env.GetBaseURL(), urlPath)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return []byte{}, err
	}

	// Add headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", apiKey)

	resp, err := sendRequest(req)
	if err != nil {
		return []byte{}, err
	}
	return resp, nil
}

func PutRequest(urlPath string, apiKey string, jsonPayload []byte) ([]byte, error) {
	return requestWithBody(urlPath, "PUT", apiKey, jsonPayload)
}

func requestWithBody(urlPath string, httpMethod string, apiKey string, jsonPayload []byte) ([]byte, error) {
	url := fmt.Sprintf("%s%s", env.GetBaseURL(), urlPath)
	req, err := http.NewRequest(httpMethod, url, bytes.NewReader(jsonPayload))
	if err != nil {
		return []byte{}, err
	}

	// Add headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", apiKey)

	resp, err := sendRequest(req)
	if err != nil {
		return []byte{}, err
	}
	return resp, nil
}

func sendRequest(req *http.Request) ([]byte, error) {
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// 400+ indicates a request error, return code / body
	fmt.Printf("resp.StatusCode :: %+v\n", resp.StatusCode)
	if 400 <= resp.StatusCode {
		return nil, fmt.Errorf("error: api error :: %d :: %q", resp.StatusCode, string(body))
	}
	return body, nil
}

// Parse response
//var response models.CreateConversationResponse
//err = json.Unmarshal(resp, &response)
//if err != nil {
//return models.CreateConversationResponse{}, err
//}

//func getRequest(urlPath string, page int) {
//	url := fmt.Sprintf("%s/%s", env.GetBaseURL(), urlPath)
//	req, err := http.NewRequest("POST", builtURL, bytes.NewReader(jsonPayload))
//	if err != nil {
//		return []Product, err
//	}
//}
