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

type registerInput struct {
	First_name  string `json:"first_name" validate:"required"`
	Middle_name string `json:"middle_name" validate:"omitempty"`
	Last_name   string `json:"last_name" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,ascii,gte=8"`
}

type RegisterUserHandler struct {
	repo          domain.UserRepository
	hasher        Hasher
	transport     *MailTransport
	keyValueStore domain.KeyValueStore
	validator     *validator.Validate
	router        *httprouter.Router
}

func NewRegisterUserHandler(
	userRepo domain.UserRepository,
	hasher Hasher,
	transport *MailTransport,
	keyValueStore domain.KeyValueStore,
	validator *validator.Validate,
	r *httprouter.Router,
) *RegisterUserHandler {
	handler := &RegisterUserHandler{
		repo:          userRepo,
		hasher:        hasher,
		transport:     transport,
		keyValueStore: keyValueStore,
		validator:     validator,
		router:        r,
	}
	handler.router.POST("/api/v1/register", handler.register())

	return handler
}

func (handler *RegisterUserHandler) register() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		input := registerInput{}
		json.NewDecoder(r.Body).Decode(&input)
		errorResponse, err := Validate(handler.validator, input)

		if err != nil {
			JSONResponse(w, errorResponse, http.StatusUnprocessableEntity)
			return
		}

		users, err := handler.repo.FindBy(r.Context(), "email", input.Email)

		if err != nil {
			JSONResponse(w, err, http.StatusBadRequest)
			return
		}

		if len(users) != 0 {
			JSONResponse(w, Error{Message: "User already exists."}, http.StatusUnprocessableEntity)
			return
		}

		hashedPassword, err := handler.hasher.Hash(input.Password)

		if err != nil {
			JSONResponse(w, err, http.StatusBadRequest)
			return
		}

		now := time.Now()
		user := domain.User{
			ID:         primitive.NewObjectID(),
			FirstName:  input.First_name,
			MiddleName: input.Middle_name,
			LastName:   input.Last_name,
			Email:      input.Email,
			Password:   hashedPassword,
			IsVerified: false,
			CreatedAt:  primitive.NewDateTimeFromTime(now),
			UpdatedAt:  primitive.NewDateTimeFromTime(now),
		}

		if err := handler.repo.Save(r.Context(), user); err != nil {
			JSONResponse(w, err, http.StatusBadRequest)
			return
		}

		randStr := GenerateRandomString(128)
		if err := SendVerificationEmail(*handler.transport, user, randStr); err != nil {
			JSONResponse(w, err, http.StatusBadRequest)
			return
		}

		tenMinutes := 10 * time.Minute
		if err := handler.keyValueStore.Store(r.Context(), randStr, user.ID.Hex(), int(tenMinutes)); err != nil {
			JSONResponse(w, err, http.StatusBadRequest)
			return
		}

		EmptyResponse(w, http.StatusCreated)
	}
}
