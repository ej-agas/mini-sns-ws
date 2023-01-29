package app

import (
	"mini-sns-ws/internal/domain"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
)

type UserHandler struct {
	repo          domain.UserRepository
	transport     *MailTransport
	keyValueStore domain.KeyValueStore
	validator     *validator.Validate
	router        *httprouter.Router
}

func NewUserHandler(
	userRepo domain.UserRepository,
	transport *MailTransport,
	keyValueStore domain.KeyValueStore,
	validator *validator.Validate,
	r *httprouter.Router,
) *UserHandler {
	h := &UserHandler{
		repo:          userRepo,
		transport:     transport,
		keyValueStore: keyValueStore,
		validator:     validator,
		router:        r,
	}

	h.routes()

	return h
}

func (h *UserHandler) routes() {
	h.router.POST("/users", h.register())
	h.router.GET("/users/verify", h.verify())
}
