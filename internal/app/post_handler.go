package app

import (
	"mini-sns-ws/internal/domain"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
)

type PostHandler struct {
	repo      domain.PostRepository
	validator *validator.Validate
	router    *httprouter.Router
}

func NewPostHandler(userRepo domain.PostRepository, validator *validator.Validate, r *httprouter.Router) *PostHandler {
	ps := &PostHandler{repo: userRepo, validator: validator, router: r}
	ps.routes()

	return ps
}

func (postHandler *PostHandler) routes() {
	postHandler.router.POST("/posts", postHandler.store())
}
