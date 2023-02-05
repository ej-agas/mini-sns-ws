package app

import (
	"context"
	"fmt"
	"mini-sns-ws/internal/domain"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthMiddleware struct {
	TokenService   JWTTokenService
	UserRepository domain.UserRepository
}

type contextKey string

const LoggedInUser contextKey = "user"

func (middleware *AuthMiddleware) Handle(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		bearerToken := r.Header.Get("Authorization")

		if bearerToken == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		plainToken := strings.Split(bearerToken, " ")[1]

		fmt.Println(plainToken)

		token, err := middleware.TokenService.FromString(plainToken)

		if err != nil {
			JSONResponse(w, Error{Message: "Invalid or expired token."}, 401)
			return
		}

		if middleware.TokenService.IsExpired(token) {
			JSONResponse(w, Error{Message: "Invalid or expired token."}, 401)
			return
		}

		fmt.Println(token.Get("user_id"))
		ctx := r.Context()
		userId, ok := token.Get("user_id")

		if !ok {
			JSONResponse(w, Error{Message: "Invalid or expired token."}, 401)
			return
		}

		str, ok := userId.(string)

		if !ok {
			JSONResponse(w, Error{Message: "Invalid or expired token."}, 401)
			return
		}

		userObjectId, err := primitive.ObjectIDFromHex(str)

		if err != nil {
			JSONResponse(w, Error{Message: "Invalid or expired token."}, 401)
			return
		}

		user, err := middleware.UserRepository.FindOneBy(ctx, "_id", userObjectId)

		if err != nil {
			JSONResponse(w, Error{Message: "User not found"}, 401)
			return
		}

		ctx = context.WithValue(ctx, LoggedInUser, user)

		next(w, r.WithContext(ctx), ps)
	}
}
