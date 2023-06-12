package domain

import (
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Following struct {
	ID        primitive.ObjectID `bson:"_id" json:"-"`
	Follower  primitive.ObjectID `bson:"follower" json:"-"`
	Following primitive.ObjectID `bson:"following" json:"following"`
	CreatedAt primitive.DateTime `bson:"created_at" json:"created_at"`
}

func (f Following) Create(followerId primitive.ObjectID, followingId primitive.ObjectID) *Following {
	now := time.Now()
	return &Following{
		ID:        primitive.NewObjectIDFromTimestamp(now),
		Follower:  followerId,
		Following: followingId,
		CreatedAt: primitive.NewDateTimeFromTime(now),
	}
}

func (f Following) Id() string {
	return f.ID.String()
}

func (f Following) String() string {
	result, err := json.Marshal(f)

	if err != nil {
		panic(err)
	}

	return string(result)
}
