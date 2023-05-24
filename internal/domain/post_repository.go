package domain

import "context"

type PostRepository interface {
	FindOne(ctx context.Context, id string) (Post, error)
	FindBy(ctx context.Context, filter Filter, findOpts FindOptions) ([]Post, error)
	FindOneBy(ctx context.Context, filter Filter) (Post, error)
	CreateFeed(ctx context.Context, ids []string, cursor string) ([]Post, error)
	Save(ctx context.Context, m Post) error
	Delete(ctx context.Context, id string) error
	DeleteBy(ctx context.Context, filter Filter) error
}
