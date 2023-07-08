package domain

import "context"

type PostRepository interface {
	FindOne(ctx context.Context, id string) (Post, error)
	FindBy(ctx context.Context, filter Filter, findOpts FindOptions) ([]Post, error)
	FindOneBy(ctx context.Context, filter Filter) (Post, error)
	CreateFeed(ctx context.Context, ids []string, page uint) ([]PostWithUser, error)
	Save(ctx context.Context, m Post) error
	Delete(ctx context.Context, id string) error
	DeleteBy(ctx context.Context, filter Filter) error
	PostsCount(ctx context.Context, user User) int64
}
