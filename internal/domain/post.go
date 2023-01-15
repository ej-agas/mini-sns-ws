package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Post struct {
	Title     string             `bson:"title" json:"title"`
	Body      string             `bson:"body" json:"body"`
	UserId    primitive.ObjectID `bson:"user_id" json:"user_id"`
	CreatedAt primitive.DateTime `bson:"created_at" json:"created_at"`
}
