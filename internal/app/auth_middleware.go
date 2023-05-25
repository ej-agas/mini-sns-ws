package app

import (
	"context"
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
			w.Header().Add("Connection", "close")
			w.WriteHeader(401)
			return
		}

		plainToken := strings.Split(bearerToken, " ")[1]
		token, err := middleware.TokenService.FromString(plainToken)

		if err != nil {
			JSONResponse(w, Error{Message: "Invalid or expired token."}, 401)
			return
		}

		if middleware.TokenService.IsExpired(token) {
			JSONResponse(w, Error{Message: "Invalid or expired token."}, 401)
			return
		}

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

		filter := domain.NewFilter()
		filter["_id"] = userObjectId

		user, err := middleware.UserRepository.FindOneBy(ctx, filter, *domain.NewFindOptions())

		if err != nil {
			JSONResponse(w, Error{Message: "User not found"}, 401)
			return
		}

		ctx = context.WithValue(ctx, LoggedInUser, user)

		next(w, r.WithContext(ctx), ps)
	}
}
