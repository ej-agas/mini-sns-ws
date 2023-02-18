package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Following struct {
	ID     primitive.ObjectID               `bson:"_id" json:"id"`
	UserId primitive.ObjectID               `bson:"user_id" json:"-"`
	List   map[primitive.ObjectID]time.Time `bson:"list" json:"list"`
}

func (following *Following) Count() int {
	return len(following.List)
}
