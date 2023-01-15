package domain

import "context"

type PostRepository interface {
	Find(ctx context.Context, id string) (Post, error)
	FindBy(ctx context.Context, field string, value interface{}) ([]Post, error)
	Save(ctx context.Context, m Post) error
	Delete(ctx context.Context, id string) error
}
