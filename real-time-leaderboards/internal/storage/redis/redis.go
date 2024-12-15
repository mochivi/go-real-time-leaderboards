package redis

import (
	"context"
	"log"

	redis "github.com/go-redis/redis/v8"
	"github.com/mochivi/go-real-time-leaderboards/conf"
)

type RedisService interface{}


type redisService struct {
	client *redis.Client
}

func NewRedisService(redisConfig conf.RedisConfig) RedisService {
	client := redis.NewClient(&redis.Options{
		Addr:     redisConfig.Addr(),
		Password: redisConfig.Password,
	})

	err := client.Ping(context.Background()).Err()
	if err != nil {
		log.Fatalf("redis connect err: %v", err)
	}

	return &redisService{client: client}
}

func (r *redisService) GetClient() *redis.Client {
	return r.client
}