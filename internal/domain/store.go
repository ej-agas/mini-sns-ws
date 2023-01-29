package domain

import "context"

type KeyValueStore interface {
	Store(ctx context.Context, key string, value interface{}, ttl int) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key ...string) error
}
