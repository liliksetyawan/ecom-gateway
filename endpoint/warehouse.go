package endpoint

import (
	"ecom-gateway/middleware"
	"ecom-gateway/service"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type WarehouseGateway struct {
	client *service.WarehouseClient
}

func NewWarehouseGateway(client *service.WarehouseClient) *WarehouseGateway {
	return &WarehouseGateway{client: client}
}

func (wg *WarehouseGateway) CreateWarehouseHandler(w http.ResponseWriter, r *http.Request) {
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	var result map[string]interface{}
	err := wg.client.CreateWarehouse(r.Header.Get("Authorization"), req, &result)
	if err != nil {
		middleware.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}

func (wg *WarehouseGateway) GetWarehouseByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid warehouse id", http.StatusBadRequest)
		return
	}

	var result map[string]interface{}
	err = wg.client.GetWarehouseByID(r.Header.Get("Authorization"), id, &result)
	if err != nil {
		middleware.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}

func (wg *WarehouseGateway) GetWarehousesHandler(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	var result map[string]interface{}
	err := wg.client.GetWarehouses(r.Header.Get("Authorization"), limit, offset, &result)
	if err != nil {
		middleware.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}
