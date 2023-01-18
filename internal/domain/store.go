package domain

import "context"

type KeyValueStore interface {
	Set(ctx context.Context, key, value string, ttlSeconds int) error
	Get(ctx context.Context, key string) (string, error)
}
