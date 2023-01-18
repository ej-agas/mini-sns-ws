package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v9"
)

type Redis struct {
	client *redis.Client
}

func (r *Redis) Store(ctx context.Context, key, value string, ttlSeconds int) error {
	return r.client.Set(ctx, key, value, time.Duration(ttlSeconds)*time.Second).Err()
}

func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()

	if err != nil {
		return "", err
	}

	return val, nil
}
