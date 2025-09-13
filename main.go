package main

import (
	"ecom-gateway/config"
	"ecom-gateway/endpoint"
	"ecom-gateway/middleware"
	"ecom-gateway/server"
	"ecom-gateway/service"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	config.LoadConfig()

	redisClient := server.NewRedisClient()

	userClient := service.NewUserClient()
	userGateway := endpoint.NewUserGateway(userClient, redisClient)

	productClient := service.NewProductClient()
	productGateway := endpoint.NewProductGateway(productClient)

	shopClient := service.NewShopClient()
	shopGateway := endpoint.NewShopGateway(shopClient)

	warehouseClient := service.NewWarehouseClient()
	warehouseGateway := endpoint.NewWarehouseGateway(warehouseClient)

	mw := middleware.NewMiddleware(redisClient)

	r := mux.NewRouter()
	r.HandleFunc("/gateway/register", userGateway.RegisterHandler).Methods("POST")
	r.HandleFunc("/gateway/login", userGateway.LoginHandler).Methods("POST")

	r.HandleFunc("/product", productGateway.CreateProductHandler).Methods("POST")
	r.HandleFunc("/product/{id:[0-9]+}", productGateway.GetProductByIDHandler).Methods("GET")
	r.HandleFunc("/product", productGateway.GetProductsHandler).Methods("GET")

	r.HandleFunc("/shop", shopGateway.CreateShopHandler).Methods("POST")
	r.HandleFunc("/shop/{id:[0-9]+}", shopGateway.GetShopByIDHandler).Methods("GET")
	r.HandleFunc("/shop", shopGateway.GetShopsHandler).Methods("GET")

	r.HandleFunc("/warehouse", warehouseGateway.CreateWarehouseHandler).Methods("POST")
	r.HandleFunc("/warehouse/{id:[0-9]+}", warehouseGateway.GetWarehouseByIDHandler).Methods("GET")
	r.HandleFunc("/warehouse", warehouseGateway.GetWarehousesHandler).Methods("GET")

	// Bungkus semua dengan middleware
	handler := mw.LoggingAndAuth(r)

	log.Println("Gateway service running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
