package mongodb

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDB struct {
	database *mongo.Database
}

func (db MongoDB) Collection(collectionName string) *mongo.Collection {
	return db.database.Collection(collectionName)
}

func NewMongoDB(databaseName string, client *mongo.Client) *MongoDB {
	return &MongoDB{client.Database(databaseName)}
}

func Connect(ctx context.Context) *mongo.Client {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGODB_URI")))

	if err != nil {
		panic(err)
	}

	if err = client.Ping(ctx, &readpref.ReadPref{}); err != nil {
		panic(err)
	}

	fmt.Println("connected to mongodb")

	return client
}
