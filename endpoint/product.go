package endpoint

import (
	"ecom-gateway/config"
	"ecom-gateway/middleware"
	"ecom-gateway/service"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type ProductGateway struct {
	client *service.ProductClient
}

func NewProductGateway(client *service.ProductClient) *ProductGateway {
	return &ProductGateway{client: client}
}

func (pg *ProductGateway) CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	ctxData, ok := r.Context().Value(middleware.ContextKey).(*middleware.ContextData)
	if !ok || ctxData == nil {
		http.Error(w, "user context not found", http.StatusUnauthorized)
		return
	}

	userID := ctxData.UserID

	w.Header().Set("Content-Type", "application/json")
	var result map[string]interface{}
	err := pg.client.CreateProduct(config.AppConfig.FixToken, userID, req, &result)
	if err != nil {
		middleware.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	json.NewEncoder(w).Encode(result)
}

func (pg *ProductGateway) GetProductByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	ctxData, ok := r.Context().Value(middleware.ContextKey).(*middleware.ContextData)
	if !ok || ctxData == nil {
		http.Error(w, "user context not found", http.StatusUnauthorized)
		return
	}

	userID := ctxData.UserID

	var result map[string]interface{}
	err = pg.client.GetProductByID(config.AppConfig.FixToken, userID, id, &result)
	if err != nil {
		middleware.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (pg *ProductGateway) GetProductsHandler(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	ctxData, ok := r.Context().Value(middleware.ContextKey).(*middleware.ContextData)
	if !ok || ctxData == nil {
		http.Error(w, "user context not found", http.StatusUnauthorized)
		return
	}

	userID := ctxData.UserID

	var result map[string]interface{}
	err := pg.client.GetProducts(config.AppConfig.FixToken, userID, limit, offset, &result)
	if err != nil {
		middleware.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
