package db

import (
	"context"
	"time"
	"url-shortener/internal/config"

	"github.com/go-redis/redis/v8"
)

//go:generate mockery --name=RedisClient --output=./mocks
type RedisClient interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
}

type redisClientImpl struct {
	client *redis.Client
}

func NewRedisClient(redisConfig *config.RedisConfig) RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr: redisConfig.Host + ":" + redisConfig.Port,
		DB:   0,
	})
	return &redisClientImpl{client: client}
}

func (r *redisClientImpl) Get(ctx context.Context, key string) *redis.StringCmd {
	return r.client.Get(ctx, key)
}

func (r *redisClientImpl) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return r.client.Set(ctx, key, value, expiration)
}

func (r *redisClientImpl) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	return r.client.Del(ctx, keys...)
}
