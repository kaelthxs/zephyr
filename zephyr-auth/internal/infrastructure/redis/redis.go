package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisClient struct {
	Client *redis.Client
	Ctx    context.Context
}

func NewRedisClient(addr, password string, db int) (*RedisClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisClient{
		Client: rdb,
		Ctx:    context.Background(),
	}, nil
}

func (r *RedisClient) Close() {
	r.Client.Close()
}

func (r *RedisClient) Ping() error {
	return r.Client.Ping(r.Ctx).Err()
}
func (r *RedisClient) Set(key string, value interface{}, expiration time.Duration) error {
	return r.Client.Set(r.Ctx, key, value, expiration).Err()
}

func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(r.Ctx, key).Result()
}

func (r *RedisClient) Del(ctx context.Context, key string) error {
	return r.Client.Del(r.Ctx, key).Err()
}

func (r *RedisClient) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return r.Client.Expire(r.Ctx, key, expiration).Err()
}

func (r *RedisClient) ExpireAt(ctx context.Context, key string, expirationTime time.Time) error {
	return r.Client.ExpireAt(r.Ctx, key, expirationTime).Err()
}
