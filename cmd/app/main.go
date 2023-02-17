package main

import (
	"context"
	"log"
	"mini-sns-ws/internal/app"
	"mini-sns-ws/internal/mongodb"
	"mini-sns-ws/internal/redis"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
)

var version string // application version

func main() {
	port := "6943"
	router := httprouter.New()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db := mongodb.NewMongoDB("sns_api", mongodb.Connect(ctx))
	validator := validator.New()

	mailCfg := app.MailTransportConfig{
		Host:     os.Getenv("MAIL_HOST"),
		Port:     os.Getenv("MAIL_PORT"),
		Password: "",
	}

	transport := app.NewMailTransport(mailCfg)
	hasher := &app.Argon2IDHasher{Params: app.Argon2IDParams{
		Memory:      64 * 1024,
		Iterations:  3,
		Parallelism: 4,
		SaltLength:  16,
		KeyLength:   32,
	}}

	userRepository := mongodb.UserRepository{UserCollection: db.Collection("users")}
	postRepository := mongodb.PostRepository{PostCollection: db.Collection("posts")}

	redis := redis.NewRedis(os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"), "", 0)
	tokenService := app.JWTTokenService{SecretKey: os.Getenv("JWT_SECRET"), Expiry: 12 * time.Hour}
	authMiddleware := app.AuthMiddleware{TokenService: tokenService, UserRepository: userRepository}

	app.NewRegisterUserHandler(userRepository, hasher, transport, redis, validator, router)
	app.NewVerifyUserHandler(userRepository, redis, router)
	app.NewLoginHandler(userRepository, hasher, tokenService, validator, router)

	app.NewCreatePostHandler(authMiddleware, postRepository, validator, router)
	app.NewGetPostsHandler(authMiddleware, postRepository, router)
	app.NewGetPostHandler(authMiddleware, postRepository, router)
	app.NewUpdatePostHandler(authMiddleware, postRepository, validator, router)
	app.NewDeletePostHandler(authMiddleware, postRepository, router)

	log.Printf("version %s listening on port %s", version, port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
