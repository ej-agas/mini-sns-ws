package app

import (
	"mini-sns-ws/internal/domain"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
)

type UserHandler struct {
	repo      domain.UserRepository
	transport domain.Transport
	validator *validator.Validate
	router    *httprouter.Router
}

func NewUserHandler(
	userRepo domain.UserRepository,
	transport domain.Transport,
	validator *validator.Validate,
	r *httprouter.Router,
) *UserHandler {
	h := &UserHandler{
		repo:      userRepo,
		transport: transport,
		validator: validator,
		router:    r,
	}

	h.routes()

	return h
}

func (h *UserHandler) routes() {
	h.router.POST("/register", h.register())
}
