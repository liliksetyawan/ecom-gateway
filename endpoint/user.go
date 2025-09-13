package endpoint

import (
	"ecom-gateway/server"
	"ecom-gateway/service"
	"encoding/json"
	"net/http"
	"strconv"
)

type UserGateway struct {
	client      *service.UserClient
	redisClient *server.RedisClient
}

func NewUserGateway(client *service.UserClient, redisClient *server.RedisClient) *UserGateway {
	return &UserGateway{
		client:      client,
		redisClient: redisClient,
	}
}

func (ug *UserGateway) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := ug.client.Register(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (ug *UserGateway) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := ug.client.Login(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	data, ok := res["data"].(map[string]interface{})
	if !ok {
		http.Error(w, "error when parsing", http.StatusBadRequest)
		return
	}

	token, ok := data["token"].(string)
	if !ok {
		http.Error(w, "error when parsing", http.StatusBadRequest)
		return
	}

	userIDFloat, ok := data["user_id"].(float64)
	if !ok {
		panic("user_id is not number")
	}
	userID := int64(userIDFloat)

	err = ug.redisClient.SetToken(token, strconv.FormatInt(userID, 10), 60*60*24)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
