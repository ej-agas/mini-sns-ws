package app

import (
	"encoding/json"
	"mini-sns-ws/internal/domain"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type unfollowUserInput struct {
	User_Id string `validate:"required" json:"user_id"`
}

type UnfollowUserHandler struct {
	authMiddleware AuthMiddleware
	validator      *validator.Validate
	repo           domain.FollowingRepository
	router         *httprouter.Router
}

func NewUnfollowUserHandler(authMiddleware AuthMiddleware, validator *validator.Validate, repo domain.FollowingRepository, router *httprouter.Router) *UnfollowUserHandler {
	handler := &UnfollowUserHandler{
		authMiddleware: authMiddleware,
		validator:      validator,
		repo:           repo,
		router:         router,
	}

	handler.router.DELETE("/api/v1/unfollow", CORS(handler.authMiddleware.Handle(handler.Handle())))

	return handler
}

func (handler UnfollowUserHandler) Handle() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		input := unfollowUserInput{}
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
		if err := handler.repo.Unfollow(r.Context(), user, userIdToFollow); err != nil {
			JSONResponse(w, Error{err.Error()}, 500)
			return
		}

		JSONResponse(w, struct {
			Message string `json:"message"`
		}{Message: "Success"}, 200)
	}
}
