package main

import (
	"context"
	"log"
	"mini-sns-ws/internal/app"
	"mini-sns-ws/internal/mongodb"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

func main() {
	port := "6942"
	router := httprouter.New()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db := mongodb.NewMongoDB("sns_api", mongodb.Connect(ctx))

	app.NewUserServer(mongodb.MongoUserRepository{UserCollection: db.Collection("users")}, router)
	app.NewPostsServer(mongodb.MongoUserRepository{UserCollection: db.Collection("posts")}, router)

	log.Printf("listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
