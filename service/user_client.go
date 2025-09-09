package service

import (
	"bytes"
	"ecom-gateway/config"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type UserClient struct {
	httpClient *http.Client
	baseURL    string
	jwtSecret  string
}

func NewUserClient() *UserClient {
	return &UserClient{
		httpClient: &http.Client{},
		baseURL:    config.AppConfig.UserServiceURL,
		jwtSecret:  config.AppConfig.JwtSecret,
	}
}

func (uc *UserClient) callAPI(method, path, token string, payload interface{}, result interface{}) error {
	// Marshal payload ke JSON
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", uc.baseURL, path), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	// Tambahkan Authorization jika ada
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := uc.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return errors.New(string(respBody))
	}

	// Decode response ke struct result
	return json.Unmarshal(respBody, result)
}

// Contoh wrapper untuk Register
func (uc *UserClient) Register(req interface{}) (map[string]interface{}, error) {
	var res map[string]interface{}
	err := uc.callAPI(http.MethodPost, "/register", "", req, &res)
	return res, err
}

// Contoh wrapper untuk Login
func (uc *UserClient) Login(req interface{}) (map[string]interface{}, error) {
	var res map[string]interface{}
	err := uc.callAPI(http.MethodPost, "/login", "", req, &res)
	return res, err
}
