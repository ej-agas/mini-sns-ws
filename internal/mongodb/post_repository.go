package mongodb

import (
	"context"
	"errors"
	"fmt"
	"mini-sns-ws/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrInvalidPostId = errors.New("mongodb: invalid post id")
)

type PostRepository struct {
	PostCollection *mongo.Collection
}

func (r PostRepository) FindOne(ctx context.Context, id string) (domain.Post, error) {
	post := domain.Post{}
	objId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return post, ErrInvalidPostId
	}

	queryError := r.PostCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&post)

	if queryError != nil {
		return post, queryError
	}

	return post, nil
}

func (r PostRepository) FindBy(ctx context.Context, filter domain.Filter, findOpts domain.FindOptions) ([]domain.Post, error) {
	var results []domain.Post

	opts := options.Find()

	sort := bson.D{}

	for field, order := range findOpts.Sort {
		sort = append(sort, bson.E{Key: field, Value: order})
	}

	opts.SetSort(sort)
	opts.SetLimit(findOpts.Limit)

	cursor, err := r.PostCollection.Find(ctx, filter, opts)

	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	fmt.Println(results)

	return results, nil
}

func (r PostRepository) FindOneBy(ctx context.Context, filter domain.Filter) (domain.Post, error) {
	var post domain.Post

	result := r.PostCollection.FindOne(ctx, filter)

	if err := result.Decode(&post); err != nil {
		return post, err
	}

	return post, nil
}

func (r PostRepository) CreateFeed(ctx context.Context, ids []string, page uint) ([]domain.PostWithUser, error) {
	perPage := 25
	skip := (page - 1) * uint(perPage)
	posts := make([]domain.PostWithUser, perPage)
	var userIds []primitive.ObjectID

	for _, v := range ids {
		userId, err := primitive.ObjectIDFromHex(v)
		if err != nil {
			return nil, fmt.Errorf("database Error: invalid id: %s", err.Error())
		}

		userIds = append(userIds, userId)
	}

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"user_id": bson.M{
					"$in": userIds,
				},
			},
		},
		{
			"$lookup": bson.M{
				"from":         "users",
				"localField":   "user_id",
				"foreignField": "_id",
				"as":           "user",
			},
		},
		{
			"$unwind": "$user",
		},
		{
			"$project": bson.M{
				"_id":        1,
				"title":      1,
				"body":       1,
				"user_id":    1,
				"created_at": 1,
				"updated_at": 1,
				"user": bson.M{
					"first_name":  1,
					"middle_name": 1,
					"last_name":   1,
					"username":    1,
					"picture":     1,
				},
			},
		},
		{
			"$skip": skip,
		},
		{
			"$limit": perPage,
		},
	}

	cursor, err := r.PostCollection.Aggregate(ctx, pipeline)

	if err != nil {
		return nil, fmt.Errorf("database Error: failed to create feed: %s", err.Error())
	}

	if err := cursor.All(ctx, &posts); err != nil {
		return nil, fmt.Errorf("database Error: failed to decode feed results: %s", err.Error())
	}

	return posts, nil
}

func (r PostRepository) Save(ctx context.Context, m domain.Post) error {
	_, err := r.PostCollection.UpdateOne(ctx, bson.D{{Key: "_id", Value: m.ID}}, bson.M{"$set": m}, options.Update().SetUpsert(true))

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

func (r PostRepository) DeleteBy(ctx context.Context, filter domain.Filter) error {
	_, err := r.PostCollection.DeleteMany(ctx, filter)

	if err != nil {
		return err
	}

	return nil
}

func (r PostRepository) PostsCount(ctx context.Context, user domain.User) int64 {
	count, err := r.PostCollection.CountDocuments(ctx, bson.M{"user_id": user.ID})

	if err != nil {
		return 0
	}

	return count
}
