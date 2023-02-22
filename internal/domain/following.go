package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Following struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Follower  primitive.ObjectID `bson:"follower" json:"follower"`
	Following primitive.ObjectID `bson:"following" json:"following"`
	CreatedAt primitive.DateTime `bson:"created_at" json:"created_at"`
}
