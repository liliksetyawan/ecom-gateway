package server

import (
	"context"
	"ecom-gateway/config"
	"github.com/go-redis/redis/v8"
	"time"
)

var ctx = context.Background()

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient() *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.AppConfig.RedisAddr,
		Password: config.AppConfig.RedisPassword,
		DB:       config.AppConfig.RedisDB,
	})
	return &RedisClient{client: rdb}
}

func (r *RedisClient) SetToken(token string, userID string, ttlSeconds int) error {
	return r.client.Set(ctx, token, userID, time.Duration(ttlSeconds)*time.Second).Err()
}

func (r *RedisClient) ValidateToken(token string) (string, error) {
	return r.client.Get(ctx, token).Result()
}
