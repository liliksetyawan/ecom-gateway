package config

import (
	"github.com/kelseyhightower/envconfig"
	"log"
)

var AppConfig *Config

type Config struct {
	UserServiceURL      string `envconfig:"USER_SERVICE_URL" required:"true"`
	ProductServiceURL   string `envconfig:"PRODUCT_SERVICE_URL" required:"true"`
	ShopServiceURL      string `envconfig:"SHOP_SERVICE_URL" required:"true"`
	WarehouseServiceURL string `envconfig:"SHOP_SERVICE_URL" required:"true"`
	OrderServiceURL     string `envconfig:"SHOP_SERVICE_URL" required:"true"`
	JwtSecret           string `envconfig:"JWT_SECRET" required:"true"`
	RedisAddr           string `envconfig:"REDIS_ADDR" default:"localhost:6379"`
	RedisPassword       string `envconfig:"REDIS_PASSWORD" default:""`
	RedisDB             int    `envconfig:"REDIS_DB" default:"0"`
}

func LoadConfig() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}
	AppConfig = &cfg
}
