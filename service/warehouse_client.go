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

type WarehouseClient struct {
	httpClient *http.Client
	baseURL    string
}

func NewWarehouseClient() *WarehouseClient {
	return &WarehouseClient{
		httpClient: &http.Client{},
		baseURL:    config.AppConfig.WarehouseServiceURL,
	}
}

func (wc *WarehouseClient) callAPI(method, path, token string, payload interface{}, result interface{}) error {
	var body io.Reader

	// encode payload kalau bukan GET
	if method != http.MethodGet && payload != nil {
		jsonBody, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		body = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", wc.baseURL, path), body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	if token != "" {
		req.Header.Set("Authorization", token)
	}

	resp, err := wc.httpClient.Do(req)
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

	if result != nil {
		return json.Unmarshal(respBody, result)
	}
	return nil
}

// CreateWarehouse : POST /warehouse
func (wc *WarehouseClient) CreateWarehouse(token string, body interface{}, result interface{}) error {
	return wc.callAPI(http.MethodPost, "/warehouse", token, body, result)
}

// GetWarehouseByID : GET /warehouse/{id}
func (wc *WarehouseClient) GetWarehouseByID(token string, id int64, result interface{}) error {
	path := fmt.Sprintf("/warehouse/%d", id)
	return wc.callAPI(http.MethodGet, path, token, nil, result)
}

// GetWarehouses : GET /warehouse?limit={}&offset={}
func (wc *WarehouseClient) GetWarehouses(token string, limit, offset int, result interface{}) error {
	path := fmt.Sprintf("/warehouse?limit=%d&offset=%d", limit, offset)
	return wc.callAPI(http.MethodGet, path, token, nil, result)
}
