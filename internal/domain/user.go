package domain

import (
	"encoding/json"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID         primitive.ObjectID `bson:"_id" json:"id"`
	FirstName  string             `bson:"first_name" json:"first_name"`
	MiddleName string             `bson:"middle_name" json:"middle_name"`
	LastName   string             `bson:"last_name" json:"last_name"`
	Email      string             `bson:"email" json:"email"`
	Password   string             `bson:"password" json:"password"`
	IsVerified bool               `bson:"is_verified" json:"is_verified"`
	VerifiedAt primitive.DateTime `bson:"verified_at" json:"verified_at"`
}

func (u User) Id() string {
	return u.ID.String()
}

func (u *User) FullName() string {
	if u.MiddleName == "" {
		return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
	}

	return fmt.Sprintf("%s %s %s", u.FirstName, u.MiddleName, u.LastName)
}

func (u User) String() string {
	result, err := json.Marshal(u)

	if err != nil {
		panic(err)
	}

	return string(result)
}
