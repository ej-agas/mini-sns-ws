package main

import (
	"context"
	"log"
	"mini-sns-ws/internal/app"
	"mini-sns-ws/internal/mongodb"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
)

func main() {
	port := "6942"
	router := httprouter.New()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db := mongodb.NewMongoDB("sns_api", mongodb.Connect(ctx))
	validator := validator.New()

	app.NewUserServer(mongodb.UserRepository{UserCollection: db.Collection("users")}, router)
	app.NewPostHandler(mongodb.PostRepository{PostCollection: db.Collection("posts")}, validator, router)

	log.Printf("listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
