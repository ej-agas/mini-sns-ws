package domain

import "context"

type PostRepository interface {
	Find(ctx context.Context, id string) (Post, error)
	FindBy(ctx context.Context, filter Filter) ([]Post, error)
	FindOneBy(ctx context.Context, filter Filter) (Post, error)
	Save(ctx context.Context, m Post) error
	Delete(ctx context.Context, id string) error
	DeleteBy(ctx context.Context, filter Filter) error
}
