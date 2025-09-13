package service

import (
	"bytes"
	"ecom-gateway/config"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type ProductClient struct {
	httpClient *http.Client
	baseURL    string
}

func NewProductClient() *ProductClient {
	return &ProductClient{
		httpClient: &http.Client{},
		baseURL:    config.AppConfig.ProductServiceURL,
	}
}

func (uc *ProductClient) callAPI(method, path, token, userID string, payload interface{}, result interface{}) error {
	var body io.Reader

	// Hanya encode payload kalau bukan GET
	if method != http.MethodGet && payload != nil {
		jsonBody, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		body = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", uc.baseURL, path), body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-UserID", userID)
	req.Header.Set("Authorization", token)

	resp, err := uc.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return errors.New(string(respBody))
	}

	return json.Unmarshal(respBody, result)
}

func (uc *ProductClient) CreateProduct(token, userID string, body interface{}, result interface{}) error {
	return uc.callAPI(http.MethodPost, "/product", userID, token, body, result)
}

func (uc *ProductClient) GetProductByID(token, userID string, id int64, result interface{}) error {
	path := fmt.Sprintf("/product/%d", id)
	return uc.callAPI(http.MethodGet, path, token, userID, nil, result)
}

func (uc *ProductClient) GetProducts(token, userID string, limit, offset int, result interface{}) error {
	path := fmt.Sprintf("/product?limit=%d&offset=%d", limit, offset)
	return uc.callAPI(http.MethodGet, path, token, userID, nil, result)
}
