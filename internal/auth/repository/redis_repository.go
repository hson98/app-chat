package auth_repository

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisRepository interface {
	SaveIDSession(c context.Context, id string, userId uuid.UUID, timeRemove time.Duration) error
	GetUserID(c context.Context, id string) (uuid.UUID, error)
	DeleteIDSession(c context.Context, id string) error
}
type authRedisRepo struct {
	redisClient *redis.Client
}

func (a authRedisRepo) DeleteIDSession(c context.Context, id string) error {
	return a.redisClient.Del(c, id).Err()
}

func (a authRedisRepo) GetUserID(c context.Context, id string) (uuid.UUID, error) {
	userid, err := a.redisClient.Get(c, id).Result()
	if err != nil {
		return [16]byte{}, err
	}
	userID, err := uuid.Parse(userid)
	if err != nil {
		return [16]byte{}, errors.New("can't convert string to uuid")
	}
	return userID, nil
}

func (a authRedisRepo) SaveIDSession(c context.Context, id string, userId uuid.UUID, timeRemove time.Duration) error {
	return a.redisClient.Set(c, id, userId, timeRemove).Err()
}

func NewAuthRedisRepo(redisClient *redis.Client) RedisRepository {
	return &authRedisRepo{redisClient: redisClient}
}
