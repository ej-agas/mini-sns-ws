package app

import (
	"encoding/json"
	"mini-sns-ws/internal/domain"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type storePostInput struct {
	Title string `validate:"required,gte=3"`
	Body  string `validate:"required,gte=5"`
}

type CreatePostHandler struct {
	authMiddleware AuthMiddleware
	repo           domain.PostRepository
	validator      *validator.Validate
	router         *httprouter.Router
}

func NewCreatePostHandler(authMiddleware AuthMiddleware, userRepo domain.PostRepository, validator *validator.Validate, router *httprouter.Router) *CreatePostHandler {
	handler := &CreatePostHandler{authMiddleware: authMiddleware, repo: userRepo, validator: validator, router: router}
	handler.router.POST("/api/v1/posts", CORS(handler.authMiddleware.Handle(handler.CreatePost())))

	return handler
}

func (handler *CreatePostHandler) CreatePost() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		input := storePostInput{}
		json.NewDecoder(r.Body).Decode(&input)
		errorResponse, err := Validate(handler.validator, input)

		if err != nil {
			JSONResponse(w, errorResponse, http.StatusUnprocessableEntity)
			return
		}

		user := (r.Context().Value(LoggedInUser)).(domain.User)
		now := time.Now()
		post := domain.Post{
			ID:        primitive.NewObjectID(),
			Title:     input.Title,
			Body:      input.Body,
			UserId:    user.ID,
			CreatedAt: primitive.NewDateTimeFromTime(now),
			UpdatedAt: primitive.NewDateTimeFromTime(now),
		}

		if err := handler.repo.Save(r.Context(), post); err != nil {
			JSONResponse(w, err, http.StatusBadRequest)
			return
		}

		response := struct {
			Id string `json:"id"`
		}{Id: post.ID.Hex()}

		JSONResponse(w, response, http.StatusCreated)
	}
}
