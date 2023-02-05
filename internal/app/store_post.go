package app

import (
	"encoding/json"
	"mini-sns-ws/internal/domain"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type storePostInput struct {
	Title string `validate:"required,gte=3"`
	Body  string `validate:"required,gte=5"`
}

func (handler *PostHandler) store() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		input := storePostInput{}
		json.NewDecoder(r.Body).Decode(&input)
		errorResponse, err := Validate(handler.validator, input)

		if err != nil {
			JSONResponse(w, errorResponse, http.StatusUnprocessableEntity)
			return
		}

		user := (r.Context().Value(LoggedInUser)).(domain.User)
		post := domain.Post{
			Title:     input.Title,
			Body:      input.Body,
			UserId:    user.ID,
			CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
		}

		if err := handler.repo.Save(r.Context(), post); err != nil {
			JSONResponse(w, err, http.StatusBadRequest)
			return
		}

		EmptyResponse(w, http.StatusCreated)
	}
}
