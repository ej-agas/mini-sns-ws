package app

import (
	"fmt"
	"mini-sns-ws/internal/domain"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FeedHandler struct {
	authMiddleware AuthMiddleware
	postRepo       domain.PostRepository
	followingRepo  domain.FollowingRepository
	router         *httprouter.Router
}

func NewFeedHandler(authMiddleware AuthMiddleware, postRepo domain.PostRepository, followingRepo domain.FollowingRepository, router *httprouter.Router) *FeedHandler {
	handler := &FeedHandler{authMiddleware: authMiddleware, postRepo: postRepo, followingRepo: followingRepo, router: router}
	handler.router.GET("/api/v1/feed", handler.authMiddleware.Handle(handler.Handle()))

	return handler
}

func (handler FeedHandler) Handle() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		user := (r.Context().Value(LoggedInUser)).(domain.User)

		var cursor string

		cursorInput := r.URL.Query().Get("cursor")
		if cursorInput != "" {
			id, err := primitive.ObjectIDFromHex(cursorInput)

			if err != nil {
				JSONResponse(w, Error{"invalid cursor"}, 400)
				return
			}

			cursor = id.Hex()
		}

		followingIds, err := handler.followingRepo.Following(r.Context(), user)

		var userIds []string
		userIds = append(userIds, user.ID.Hex())

		for _, following := range followingIds {
			userIds = append(userIds, following.ID.Hex())
		}

		if err != nil {
			JSONResponse(w, Error{fmt.Sprintf("Error getting list of following: %s", err.Error())}, 500)
			return
		}

		feed, err := handler.postRepo.CreateFeed(r.Context(), userIds, cursor)

		if err != nil {
			JSONResponse(w, Error{fmt.Sprintf("Error getting feed: %s", err.Error())}, 500)
			return
		}

		JSONResponse(w, feed, 200)
	}
}
