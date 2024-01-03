package database

import (
	"github.com/redis/go-redis/v9"
	"log"
	"sync"
)

var (
	redisClient *redis.Client
	once        sync.Once
)

func GetRedisClient(redisHost string) *redis.Client {
	once.Do(func() {
		opts, err := redis.ParseURL(redisHost)
		if err != nil {
			log.Printf("error connect redis %v", err)
			panic(err)
		}
		redisClient = redis.NewClient(opts)
	})
	return redisClient
}
