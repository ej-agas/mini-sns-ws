package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FollowingRepository interface {
	Follow(ctx context.Context, follower User, userToFollow primitive.ObjectID) error
	Unfollow(ctx context.Context, follower User, userToUnfollow primitive.ObjectID) error
}
