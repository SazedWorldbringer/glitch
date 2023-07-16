package database

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

func CreateClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: os.Getenv("DB_ADDRESS"),
		Password: os.Getenv("DB_PASSWORD"),
		DB: 0,
	})

	return client
}
