package domain

import (
	"encoding/json"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID         primitive.ObjectID `bson:"_id" json:"id"`
	FirstName  string             `bson:"first_name" json:"first_name"`
	MiddleName string             `bson:"middle_name" json:"middle_name"`
	LastName   string             `bson:"last_name" json:"last_name"`
	Bio        string             `bson:"bio" json:"bio"`
	Email      string             `bson:"email" json:"email"`
	Password   string             `bson:"password" json:"-"`
	Picture    string             `bson:"picture" json:"picture"`
	IsVerified bool               `bson:"is_verified" json:"is_verified"`
	VerifiedAt primitive.DateTime `bson:"verified_at,omitempty" json:"verified_at,omitempty"`
	CreatedAt  primitive.DateTime `bson:"created_at" json:"created_at"`
	UpdatedAt  primitive.DateTime `bson:"updated_at" json:"updated_at"`
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

func (u *User) Verify() {
	if u.IsVerified {
		return
	}

	now := time.Now()
	u.IsVerified = true
	u.VerifiedAt = primitive.NewDateTimeFromTime(now)
	u.UpdatedAt = primitive.NewDateTimeFromTime(now)
}

func (u User) String() string {
	result, err := json.Marshal(u)

	if err != nil {
		panic(err)
	}

	return string(result)
}
