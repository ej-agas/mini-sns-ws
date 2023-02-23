package app

import (
	"mini-sns-ws/internal/domain"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type FollowUserHandler struct {
	authMiddleware AuthMiddleware
	repo           domain.FollowingRepository
	router         *httprouter.Router
}

func NewFollowUserHandler(authMiddleware AuthMiddleware, repo domain.FollowingRepository, router *httprouter.Router) *FollowUserHandler {
	return &FollowUserHandler{
		authMiddleware: authMiddleware,
		repo:           repo,
		router:         router,
	}
}

func (handler FollowUserHandler) Handle() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	}
}
