package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Followers struct {
	ID     primitive.ObjectID               `bson:"_id" json:"id"`
	UserId primitive.ObjectID               `bson:"user_id" json:"-"`
	List   map[primitive.ObjectID]time.Time `bson:"list" json:"list"`
}

func (followers *Followers) Count() int {
	return len(followers.List)
}
