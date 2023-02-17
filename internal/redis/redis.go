package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v9"
)

type Redis struct {
	client *redis.Client
}

func (r *Redis) Store(ctx context.Context, key string, value interface{}, ttl int) error {
	return r.client.Set(ctx, key, value, time.Duration(ttl)*time.Second).Err()
}

func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()

	if err != nil {
		return "", err
	}

	return val, nil
}

func (r *Redis) Delete(ctx context.Context, keys ...string) error {
	_, err := r.client.Del(ctx, keys...).Result()

	if err != nil {
		return err
	}

	return nil
}

func NewRedis(host string, port string, password string, db int) *Redis {
	return &Redis{
		redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", host, port),
			Password: password,
			DB:       db,
		}),
	}
}
