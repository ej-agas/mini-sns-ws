package app

import (
	"fmt"
	"mini-sns-ws/internal/domain"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type FeedHandler struct {
	authMiddleware AuthMiddleware
	postRepo       domain.PostRepository
	followingRepo  domain.FollowingRepository
	router         *httprouter.Router
}

type FeedData struct {
	Data []domain.PostWithUser `json:"data"`
}

func NewFeedHandler(authMiddleware AuthMiddleware, postRepo domain.PostRepository, followingRepo domain.FollowingRepository, router *httprouter.Router) *FeedHandler {
	handler := &FeedHandler{authMiddleware: authMiddleware, postRepo: postRepo, followingRepo: followingRepo, router: router}
	handler.router.GET("/api/v1/feed", CORS(handler.authMiddleware.Handle(handler.Handle())))

	return handler
}

func (handler FeedHandler) Handle() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		user := (r.Context().Value(LoggedInUser)).(domain.User)

		page, err := strconv.ParseUint(r.URL.Query().Get("page"), 10, 64)
		if err != nil {
			JSONResponse(w, Error{"invalid page."}, 400)
			return
		}

		followingIds, err := handler.followingRepo.Following(r.Context(), user)

		var userIds []string
		userIds = append(userIds, user.ID.Hex())

		for _, following := range followingIds {
			userIds = append(userIds, following.Following.Hex())
		}

		if err != nil {
			JSONResponse(w, Error{fmt.Sprintf("Error getting list of following: %s", err.Error())}, 500)
			return
		}

		feed, err := handler.postRepo.CreateFeed(r.Context(), userIds, uint(page))

		if err != nil {
			JSONResponse(w, Error{fmt.Sprintf("Error getting feed: %s", err.Error())}, 500)
			return
		}

		data := FeedData{feed}

		JSONResponse(w, data, 200)
	}
}
