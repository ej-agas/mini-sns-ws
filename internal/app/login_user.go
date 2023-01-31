package app

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type loginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,ascii"`
}

func (h *UserHandler) login() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		input := loginInput{}
		json.NewDecoder(r.Body).Decode(&input)
		errorResponse, err := Validate(h.validator, input)

		if err != nil {
			JSONResponse(w, errorResponse, http.StatusUnprocessableEntity)
			return
		}

		user, err := h.repo.FindOneBy(r.Context(), "email", input.Email)

		if err != nil {
			JSONResponse(w, Error{Message: "User not found."}, 400)
			return
		}

		ok, err := h.hasher.Check(input.Password, user.Password)

		if err != nil {
			JSONResponse(w, err.Error(), 400)
			return
		}

		if !ok {
			JSONResponse(w, Error{Message: "Invalid password."}, 422)
			return
		}

		JSONResponse(w, user, 200)
	}
}
