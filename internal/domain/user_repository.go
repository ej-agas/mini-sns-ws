package domain

import "context"

type UserRepository interface {
	Find(ctx context.Context, id string) (User, error)
	FindBy(ctx context.Context, filter Filter, findOpts FindOptions) ([]User, error)
	FindOneBy(ctx context.Context, filter Filter, findOpts FindOptions) (User, error)
	Search(ctx context.Context, query string) ([]User, error)
	Save(ctx context.Context, m User) error
	Delete(ctx context.Context, id string) error
}
