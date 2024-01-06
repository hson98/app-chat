package middlewares

import (
	"github.com/redis/go-redis/v9"
	"hson98/app-chat/pkg/myjwt"
)

type MiddlewareManager struct {
	origins     []string
	jwtMaker    myjwt.Maker
	clientRedis *redis.Client
}

// contructor
func NewMiddlewareManager(origins []string, jwtMaker myjwt.Maker, clientRedis *redis.Client) *MiddlewareManager {
	return &MiddlewareManager{origins: origins, jwtMaker: jwtMaker, clientRedis: clientRedis}
}
