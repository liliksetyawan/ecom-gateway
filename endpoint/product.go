package endpoint

import (
	"ecom-gateway/service"
	"encoding/json"
	"net/http"
)

type ProductGateway struct {
	client *service.ProductClient
}

func NewProductGateway(client *service.ProductClient) *ProductGateway {
	return &ProductGateway{client: client}
}

// Contoh Get List Product (proxy ke product-service)
func (pg *ProductGateway) GetProductsHandler(w http.ResponseWriter, r *http.Request) {
	var result map[string]interface{}

	err := pg.client.GetProduct(r.Header.Get("Authorization"), &result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
