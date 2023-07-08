package domain

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Title     string             `bson:"title" json:"title"`
	Body      string             `bson:"body" json:"body"`
	UserId    primitive.ObjectID `bson:"user_id" json:"user_id"`
	CreatedAt primitive.DateTime `bson:"created_at" json:"created_at"`
	UpdatedAt primitive.DateTime `bson:"updated_at" json:"updated_at"`
}

func (p Post) Id() string {
	return p.ID.String()
}

func (p Post) String() string {
	result, err := json.Marshal(p)

	if err != nil {
		panic(err)
	}

	return string(result)
}

type PostWithUser struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Title     string             `bson:"title" json:"title"`
	Body      string             `bson:"body" json:"body"`
	UserId    primitive.ObjectID `bson:"user_id" json:"user_id"`
	CreatedAt primitive.DateTime `bson:"created_at" json:"created_at"`
	UpdatedAt primitive.DateTime `bson:"updated_at" json:"updated_at"`
	User      User               `bson:"user" json:"user"`
}
