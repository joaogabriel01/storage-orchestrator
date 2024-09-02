package main

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type RedisStorageUnit struct {
	client *redis.Client
}

func NewRedisStorageUnit() *RedisStorageUnit {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	return &RedisStorageUnit{client: redisClient}
}

func (r *RedisStorageUnit) Save(ctx context.Context, key string, value string) error {
	return r.client.Set(ctx, key, value, 0).Err()
}

func (r *RedisStorageUnit) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

func (r *RedisStorageUnit) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
