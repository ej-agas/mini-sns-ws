package main

import (
	"context"
	"log"
	"mini-sns-ws/internal/app"
	"mini-sns-ws/internal/mongodb"
	"mini-sns-ws/internal/redis"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
)

func main() {
	port := "6943"
	router := httprouter.New()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db := mongodb.NewMongoDB("sns_api", mongodb.Connect(ctx))
	validator := validator.New()

	mailCfg := app.MailTransportConfig{
		Host:     "0.0.0.0",
		Port:     "32769",
		Password: "",
	}

	transport := app.NewMailTransport(mailCfg)

	redis := redis.NewRedis("0.0.0.0:7000", "", 0)

	app.NewUserHandler(mongodb.UserRepository{UserCollection: db.Collection("users")}, transport, redis, validator, router)
	app.NewPostHandler(mongodb.PostRepository{PostCollection: db.Collection("posts")}, validator, router)

	log.Printf("listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
