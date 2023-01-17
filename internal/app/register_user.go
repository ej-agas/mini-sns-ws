package app

import (
	"encoding/json"
	"mini-sns-ws/internal/domain"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type registerInput struct {
	First_name  string `json:"first_name" validate:"required,alpha"`
	Middle_name string `json:"middle_name" validate:"omitempty,alpha"`
	Last_name   string `json:"last_name" validate:"required,alpha"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,alphanum"`
}

func (h *UserHandler) register() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		input := registerInput{}
		json.NewDecoder(r.Body).Decode(&input)
		errorResponse, err := Validate(h.validator, input)

		if err != nil {
			JSONResponse(w, errorResponse, http.StatusUnprocessableEntity)
			return
		}

		user := domain.User{
			ID:         primitive.NewObjectID(),
			FirstName:  input.First_name,
			MiddleName: input.Middle_name,
			LastName:   input.Last_name,
			Email:      input.Email,
			Password:   input.Password,
			IsVerified: false,
		}

		if err := h.repo.Save(r.Context(), user); err != nil {
			JSONResponse(w, err, http.StatusBadRequest)
			return
		}

		EmptyResponse(w, http.StatusCreated)
	}
}
