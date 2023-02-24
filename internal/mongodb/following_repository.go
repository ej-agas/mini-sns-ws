package mongodb

import (
	"context"
	"errors"
	"mini-sns-ws/internal/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrAlreadyFollowingUser       = errors.New("mongodb: aready following user")
	ErrUserToFollowDoesNotExist   = errors.New("mongodb: user to follow does not exist")
	ErrUserToUnfollowDoesNotExist = errors.New("mongodb: user to unfollow does not exist")
)

type FollowingRepository struct {
	FollowingCollection *mongo.Collection
	UserCollection      *mongo.Collection
}

func (repository FollowingRepository) Follow(ctx context.Context, follower domain.User, userToFollow primitive.ObjectID) error {
	var following domain.Following
	var user domain.User

	userResult := repository.UserCollection.FindOne(ctx, bson.M{"_id": userToFollow})
	if err := userResult.Decode(&user); err != nil {
		return ErrUserToFollowDoesNotExist
	}

	result := repository.FollowingCollection.FindOne(ctx, bson.M{"follower": follower.ID, "following": userToFollow})
	if err := result.Decode(&following); err == nil {
		return ErrAlreadyFollowingUser
	}

	following.ID = primitive.NewObjectID()
	following.Follower = follower.ID
	following.Following = userToFollow
	following.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	if _, err := repository.FollowingCollection.InsertOne(ctx, following); err != nil {
		return err
	}

	return nil
}

func (repository FollowingRepository) Unfollow(ctx context.Context, follower domain.User, userToUnfollow primitive.ObjectID) error {
	if _, err := repository.FollowingCollection.DeleteOne(ctx, bson.M{"follower": follower.ID, "following": userToUnfollow}); err != nil {
		return err
	}

	return nil
}
