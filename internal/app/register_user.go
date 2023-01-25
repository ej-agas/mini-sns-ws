package app

import (
	"encoding/json"
	"fmt"
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

		users, err := h.repo.FindBy(r.Context(), "email", input.Email)

		if err != nil {
			JSONResponse(w, err, http.StatusBadRequest)
			return
		}

		if len(users) != 0 {
			JSONResponse(w, Error{Message: "User already exists."}, http.StatusUnprocessableEntity)
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

		// if err := h.repo.Save(r.Context(), user); err != nil {
		// 	JSONResponse(w, err, http.StatusBadRequest)
		// 	return
		// }

		randStr := GenerateRandomString(128)

		data := struct {
			Name string
			URL  string
		}{
			Name: user.FullName(),
			URL:  "http://localhost:6943/users/verify?token=" + randStr,
		}

		mail, err := NewMail("internal/templates/VerifyAccount.html", []string{"foo@email.com"}, "noreply@mini-sns.com", "Verify", data)

		if err != nil {
			fmt.Println(err)
			JSONResponse(w, err, http.StatusBadRequest)
			return
		}

		if err := h.transport.Send(mail); err != nil {
			JSONResponse(w, err, http.StatusBadRequest)
			fmt.Println(err)
			return
		}

		EmptyResponse(w, http.StatusCreated)
	}
}
