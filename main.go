package main

import (
	"ecom-gateway/config"
	"ecom-gateway/endpoint"
	"ecom-gateway/middleware"
	"ecom-gateway/server"
	"ecom-gateway/service"
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
	userGateway := endpoint.NewUserGateway(userClient)

	productClient := service.NewProductClient()
	productGateway := endpoint.NewProductGateway(productClient)

	mw := middleware.NewMiddleware(redisClient)

	mux := http.NewServeMux()
	mux.HandleFunc("/gateway/register", userGateway.RegisterHandler)
	mux.HandleFunc("/gateway/login", userGateway.LoginHandler)
	mux.HandleFunc("/gateway/products", productGateway.GetProductsHandler)

	// Bungkus semua dengan middleware
	handler := mw.LoggingAndAuth(mux)

	log.Println("Gateway service running on port 8081...")
	log.Fatal(http.ListenAndServe(":8081", handler))
}
