package app

import (
	"encoding/json"
	"errors"
	"mini-sns-ws/internal/domain"
	"mini-sns-ws/internal/mongodb"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrInvalidUserId = errors.New("invalid user id")
)

type followUserInput struct {
	User_Id string `validate:"required" json:"user_id"`
}

type FollowUserHandler struct {
	authMiddleware AuthMiddleware
	validator      *validator.Validate
	repo           domain.FollowingRepository
	router         *httprouter.Router
}

func NewFollowUserHandler(authMiddleware AuthMiddleware, validator *validator.Validate, repo domain.FollowingRepository, router *httprouter.Router) *FollowUserHandler {
	handler := &FollowUserHandler{
		authMiddleware: authMiddleware,
		validator:      validator,
		repo:           repo,
		router:         router,
	}

	handler.router.POST("/api/v1/follow", CORS(handler.authMiddleware.Handle(handler.Handle())))

	return handler
}

func (handler FollowUserHandler) Handle() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		input := followUserInput{}
		json.NewDecoder(r.Body).Decode(&input)
		errorResponse, err := Validate(handler.validator, input)

		if err != nil {
			JSONResponse(w, errorResponse, http.StatusUnprocessableEntity)
			return
		}

		userIdToFollow, err := primitive.ObjectIDFromHex(input.User_Id)

		if err != nil {
			JSONResponse(w, Error{ErrInvalidUserId.Error()}, 422)
			return
		}

		user := (r.Context().Value(LoggedInUser)).(domain.User)
		err = handler.repo.Follow(r.Context(), user, userIdToFollow)

		if err == mongodb.ErrAlreadyFollowingUser {
			JSONResponse(w, Error{"You are already following this user"}, 422)
			return
		}

		if err != nil {
			JSONResponse(w, Error{err.Error()}, 500)
			return
		}

		JSONResponse(w, struct {
			Message string `json:"message"`
		}{Message: "Success"}, 200)
	}
}
