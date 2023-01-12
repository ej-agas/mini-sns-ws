package domain

import "context"

type UserRepository interface {
	Find(ctx context.Context, id string) (User, error)
	FindBy(ctx context.Context, field string, value interface{}) ([]User, error)
	Save(ctx context.Context, m User) error
	Delete(ctx context.Context, id string) error
}
