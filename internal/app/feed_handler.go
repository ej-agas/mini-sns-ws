package app

import (
	"mini-sns-ws/internal/domain"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type FeedHandler struct {
	authMiddlerware AuthMiddleware
	postRepo        domain.PostRepository
	followingRepo   domain.FollowingRepository
	router          *httprouter.Router
}

func NewFeedHandler(authMiddleware AuthMiddleware, postRepo domain.PostRepository, followingRepo domain.FollowingRepository, router *httprouter.Router) *FeedHandler {
	handler := &FeedHandler{authMiddlerware: authMiddleware, postRepo: postRepo, followingRepo: followingRepo, router: router}
	handler.router.GET("/api/v1/feed", handler.authMiddlerware.Handle(handler.Handle()))

	return handler
}

func (handler FeedHandler) Handle() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	}
}
