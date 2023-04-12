package app

import (
	"encoding/json"
	"mini-sns-ws/internal/domain"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
)

type LoginHandler struct {
	repo         domain.UserRepository
	hasher       Hasher
	tokenService TokenService
	validator    *validator.Validate
	router       *httprouter.Router
}

func NewLoginHandler(
	userRepo domain.UserRepository,
	hasher Hasher,
	tokenService TokenService,
	validator *validator.Validate,
	r *httprouter.Router,
) *LoginHandler {
	handler := &LoginHandler{
		repo:         userRepo,
		hasher:       hasher,
		tokenService: tokenService,
		validator:    validator,
		router:       r,
	}
	handler.router.POST("/api/v1/login", CORS(handler.login()))
	return handler
}

type loginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,ascii"`
}

type loginResponse struct {
	Message string `json:"token"`
}

func (handler *LoginHandler) login() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		input := loginInput{}
		json.NewDecoder(r.Body).Decode(&input)
		errorResponse, err := Validate(handler.validator, input)

		if err != nil {
			JSONResponse(w, errorResponse, http.StatusUnprocessableEntity)
			return
		}

		user, err := handler.repo.FindOneBy(r.Context(), "email", input.Email)

		if err != nil {
			JSONResponse(w, Error{Message: "User not found."}, 400)
			return
		}

		ok, err := handler.hasher.Check(input.Password, user.Password)

		if err != nil {
			JSONResponse(w, err.Error(), 400)
			return
		}

		if !ok {
			JSONResponse(w, Error{Message: "Invalid password."}, 422)
			return
		}

		jwtToken, err := handler.tokenService.GenerateFor(user)

		if err != nil {
			JSONResponse(w, Error{Message: err.Error()}, 400)
			return
		}

		JSONResponse(w, loginResponse{Message: jwtToken}, 200)
	}
}
