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

type ShopClient struct {
	httpClient *http.Client
	baseURL    string
}

func NewShopClient() *ShopClient {
	return &ShopClient{
		httpClient: &http.Client{},
		baseURL:    config.AppConfig.ShopServiceURL,
	}
}

func (sc *ShopClient) callAPI(method, path, token string, payload interface{}, result interface{}) error {
	var body io.Reader

	// encode payload kalau bukan GET
	if method != http.MethodGet && payload != nil {
		jsonBody, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		body = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", sc.baseURL, path), body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	if token != "" {
		req.Header.Set("Authorization", token)
	}

	resp, err := sc.httpClient.Do(req)
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

// CreateShop : POST /shop
func (sc *ShopClient) CreateShop(token string, body interface{}, result interface{}) error {
	return sc.callAPI(http.MethodPost, "/shop", token, body, result)
}

// GetShopByID : GET /shop/{id}
func (sc *ShopClient) GetShopByID(token string, id int64, result interface{}) error {
	path := fmt.Sprintf("/shop/%d", id)
	return sc.callAPI(http.MethodGet, path, token, nil, result)
}

// GetShops : GET /shop?limit={}&offset={}
func (sc *ShopClient) GetShops(token string, limit, offset int, result interface{}) error {
	path := fmt.Sprintf("/shop?limit=%d&offset=%d", limit, offset)
	return sc.callAPI(http.MethodGet, path, token, nil, result)
}
