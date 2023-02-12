package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Post struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Title     string             `bson:"title" json:"title"`
	Body      string             `bson:"body" json:"body"`
	UserId    primitive.ObjectID `bson:"user_id" json:"-"`
	CreatedAt primitive.DateTime `bson:"created_at" json:"created_at"`
	UpdatedAt primitive.DateTime `bson:"updated_at" json:"updated_at"`
}
