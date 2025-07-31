package cache

import (
    "context"
    "time"

    "github.com/redis/go-redis/v9"
)

type RedisClient struct {
    Client *redis.Client
    Ctx    context.Context
}

func NewRedisClient(addr, password string, db int) *RedisClient {
    rdb := redis.NewClient(&redis.Options{
        Addr:     addr,
        Password: password,
        DB:       db,
    })

    return &RedisClient{
        Client: rdb,
        Ctx:    context.Background(),
    }
}

func (r *RedisClient) Set(key, value string, expiration time.Duration) error {
    return r.Client.Set(r.Ctx, key, value, expiration).Err()
}

func (r *RedisClient) Get(key string) (string, error) {
    return r.Client.Get(r.Ctx, key).Result()
}

func (r *RedisClient) Delete(key string) error {
    return r.Client.Del(r.Ctx, key).Err()
}

func (r *RedisClient) SetEmailCode(email, code string, expiration time.Duration) error {
    if err := r.Client.Set(r.Ctx, "email_code:"+email, code, expiration).Err(); err != nil {
        return err
    }
    return r.Client.Set(r.Ctx, "email_code_reverse:"+code, email, expiration).Err()
}

func (r *RedisClient) GetEmailByCode(code string) (string, error) {
    return r.Client.Get(r.Ctx, "email_code_reverse:"+code).Result()
}

func (r *RedisClient) DeleteEmailCode(email string) error {
    code, err := r.Client.Get(r.Ctx, "email_code:"+email).Result()
    if err != nil {
        if err == redis.Nil {
            return nil
        }
        return err
    }

    err1 := r.Client.Del(r.Ctx, "email_code:"+email).Err()
    err2 := r.Client.Del(r.Ctx, "email_code_reverse:"+code).Err()
    if err1 != nil {
        return err1
    }
    return err2
}


