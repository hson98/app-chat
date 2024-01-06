package database

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"sync"
)

var (
	redisClient *redis.Client
	once        sync.Once
)

func GetRedisClient(redisHost string, pass string) *redis.Client {
	once.Do(func() {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     redisHost,
			Password: "redis", // no password set
			DB:       0,       // use default DB
		})
		_, err := redisClient.Ping(context.Background()).Result()
		if err != nil {
			log.Fatalf("Lỗi connect redis%s", err.Error())
			panic(err)
		}
	})
	return redisClient
}
