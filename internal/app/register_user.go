package app

import (
	"encoding/json"
	"fmt"
	"mini-sns-ws/internal/domain"
	"net/http"
	"time"

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

		now := time.Now()
		user := domain.User{
			ID:         primitive.NewObjectID(),
			FirstName:  input.First_name,
			MiddleName: input.Middle_name,
			LastName:   input.Last_name,
			Email:      input.Email,
			Password:   input.Password,
			IsVerified: false,
			CreatedAt:  primitive.NewDateTimeFromTime(now),
			UpdatedAt:  primitive.NewDateTimeFromTime(now),
		}

		if err := h.repo.Save(r.Context(), user); err != nil {
			fmt.Println(err)
			JSONResponse(w, err, http.StatusBadRequest)
			return
		}

		randStr := GenerateRandomString(128)

		if err := sendVerificationEmail(*h.transport, user, randStr); err != nil {
			JSONResponse(w, err, http.StatusBadRequest)
			return
		}

		tenMinutes := 600
		if err := h.keyValueStore.Store(r.Context(), randStr, user.ID.Hex(), tenMinutes); err != nil {
			JSONResponse(w, err, http.StatusBadRequest)
			return
		}

		EmptyResponse(w, http.StatusCreated)
	}
}

func sendVerificationEmail(transport MailTransport, user domain.User, verificationToken string) error {
	data := struct {
		Name string
		URL  string
	}{
		Name: user.FullName(),
		URL:  "http://localhost:6943/users/verify?token=" + verificationToken,
	}

	mailTo := []string{"foo@email.com"}
	subject := "Verify your account"

	mail, err := NewMail("internal/templates/VerifyAccount.html", mailTo, "noreply@mini-sns.com", subject, data)

	if err != nil {
		return err
	}

	if err := transport.Send(mail); err != nil {
		return err
	}

	return nil
}
