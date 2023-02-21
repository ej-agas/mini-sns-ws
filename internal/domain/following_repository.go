package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type FollowingRepository interface {
	Follow(follower User, userToFollow primitive.ObjectID) error
	Unfollow(follower User, userToUnfollow primitive.ObjectID) error
}
