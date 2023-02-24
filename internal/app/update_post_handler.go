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

type updatePostInput struct {
	Title string `validate:"omitempty,gte=3"`
	Body  string `validate:"omitempty,gte=5"`
}

type UpdatePostHandler struct {
	authMiddleware AuthMiddleware
	repo           domain.PostRepository
	validator      *validator.Validate
	router         *httprouter.Router
}

func NewUpdatePostHandler(authMiddleware AuthMiddleware, repo domain.PostRepository, validator *validator.Validate, router *httprouter.Router) *UpdatePostHandler {
	handler := &UpdatePostHandler{authMiddleware: authMiddleware, repo: repo, validator: validator, router: router}
	handler.router.PUT("/api/v1/posts/:id", handler.authMiddleware.Handle(handler.updatePost()))

	return handler
}

func (handler UpdatePostHandler) updatePost() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		input := updatePostInput{}
		json.NewDecoder(r.Body).Decode(&input)
		errorResponse, err := Validate(handler.validator, input)

		if err != nil {
			JSONResponse(w, errorResponse, 422)
			return
		}

		user := r.Context().Value(LoggedInUser).(domain.User)
		postId, err := primitive.ObjectIDFromHex(ps.ByName("id"))

		if err != nil {
			JSONResponse(w, Error{Message: "invalid id"}, 404)
			return
		}

		filter := domain.NewFilter()
		filter["_id"] = postId
		filter["user_id"] = user.ID

		post, err := handler.repo.FindOneBy(r.Context(), filter)

		if err != nil {
			JSONResponse(w, Error{Message: "post not found."}, 404)
			return
		}

		if input.Title != "" {
			post.Title = input.Title
		}

		if input.Body != "" {
			post.Body = input.Body
		}

		post.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

		if err := handler.repo.Save(r.Context(), post); err != nil {
			JSONResponse(w, Error{Message: err.Error()}, 400)
			return
		}

		JSONResponse(w, post, 200)
	}
}
