package mongodb

import (
	"context"
	"errors"
	"mini-sns-ws/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrInvalidUserId = errors.New("mongodb: invalid user id")
)

type UserRepository struct {
	UserCollection *mongo.Collection
}

func (r UserRepository) Find(ctx context.Context, id string) (domain.User, error) {
	user := domain.User{}
	objId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return user, ErrInvalidUserId
	}

	queryError := r.UserCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)

	if queryError != nil {
		return user, queryError
	}

	return user, nil
}

func (r UserRepository) FindBy(ctx context.Context, field string, value interface{}) ([]domain.User, error) {
	var results []domain.User

	cursor, err := r.UserCollection.Find(ctx, bson.D{{Key: field, Value: value}})

	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (r UserRepository) FindOneBy(ctx context.Context, field string, value interface{}) (domain.User, error) {
	var user domain.User

	result := r.UserCollection.FindOne(ctx, bson.D{{Key: field, Value: value}})

	if err := result.Decode(&user); err != nil {
		return user, err
	}

	return user, nil
}

func (r UserRepository) Search(ctx context.Context, query string) ([]domain.User, error) {
	var results []domain.User

	filter := bson.M{
		"$text": []bson.M{
			{"$search": query},
		},
	}

	cursor, err := r.UserCollection.Find(ctx, filter)

	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (r UserRepository) Save(ctx context.Context, m domain.User) error {
	_, err := r.UserCollection.UpdateOne(ctx, bson.D{{Key: "_id", Value: m.ID}}, bson.M{"$set": m}, options.Update().SetUpsert(true))

	if err != nil {
		return err
	}

	return nil
}

func (r UserRepository) Delete(ctx context.Context, id string) error {
	_, err := r.UserCollection.DeleteOne(ctx, bson.D{{Key: "_id", Value: id}})

	if err != nil {
		return err
	}

	return nil
}
