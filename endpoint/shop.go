package endpoint

import (
	"ecom-gateway/middleware"
	"ecom-gateway/service"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type ShopGateway struct {
	client *service.ShopClient
}

func NewShopGateway(client *service.ShopClient) *ShopGateway {
	return &ShopGateway{client: client}
}

// CreateShopHandler : POST /shop
func (sg *ShopGateway) CreateShopHandler(w http.ResponseWriter, r *http.Request) {
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	var result map[string]interface{}
	err := sg.client.CreateShop(r.Header.Get("Authorization"), req, &result)
	if err != nil {
		middleware.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}

// GetShopByIDHandler : GET /shop/{id}
func (sg *ShopGateway) GetShopByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid shop id", http.StatusBadRequest)
		return
	}

	var result map[string]interface{}
	err = sg.client.GetShopByID(r.Header.Get("Authorization"), id, &result)
	if err != nil {
		middleware.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}

// GetShopsHandler : GET /shop?limit=&offset=&search=
func (sg *ShopGateway) GetShopsHandler(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	var result map[string]interface{}
	err := sg.client.GetShops(r.Header.Get("Authorization"), limit, offset, &result)
	if err != nil {
		middleware.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}
