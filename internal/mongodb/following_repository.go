package mongodb

import (
	"context"
	"errors"
	"mini-sns-ws/internal/domain"

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

	if _, err := repository.FollowingCollection.InsertOne(ctx, following.Create(follower.ID, userToFollow)); err != nil {
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

func (repository FollowingRepository) Following(ctx context.Context, user domain.User) ([]domain.Following, error) {
	var results []domain.Following

	cursor, err := repository.FollowingCollection.Find(ctx, bson.M{"follower": user.ID})

	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (repository FollowingRepository) Followers(ctx context.Context, user domain.User) ([]domain.Following, error) {
	var results []domain.Following

	cursor, err := repository.FollowingCollection.Find(ctx, bson.M{"following": user.ID})

	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
