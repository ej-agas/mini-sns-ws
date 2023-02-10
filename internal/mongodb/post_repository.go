package mongodb

import (
	"context"
	"errors"
	"mini-sns-ws/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostRepository struct {
	PostCollection *mongo.Collection
}

func (r PostRepository) Find(ctx context.Context, id string) (domain.Post, error) {
	post := domain.Post{}
	objId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return post, errors.New("invalid ObjectID")
	}

	queryError := r.PostCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&post)

	if queryError != nil {
		return post, queryError
	}

	return post, nil
}

func (r PostRepository) FindBy(ctx context.Context, field string, value interface{}) ([]domain.Post, error) {
	var results []domain.Post

	cursor, err := r.PostCollection.Find(ctx, bson.D{{Key: field, Value: value}})

	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (r PostRepository) FindOneBy(ctx context.Context, field string, value interface{}) (domain.Post, error) {
	var post domain.Post

	result := r.PostCollection.FindOne(ctx, bson.D{{Key: field, Value: value}})

	if err := result.Decode(&post); err != nil {
		return post, err
	}

	return post, nil
}

func (r PostRepository) Save(ctx context.Context, m domain.Post) error {
	_, err := r.PostCollection.InsertOne(ctx, m)

	if err != nil {
		return err
	}

	return nil
}

func (r PostRepository) Delete(ctx context.Context, id string) error {
	_, err := r.PostCollection.DeleteOne(ctx, bson.D{{Key: "_id", Value: id}})

	if err != nil {
		return err
	}

	return nil
}
