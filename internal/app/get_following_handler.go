package app

import (
	"mini-sns-ws/internal/domain"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type GetFollowingHandler struct {
	authMiddlerware AuthMiddleware
	followingRepo   domain.FollowingRepository
	router          *httprouter.Router
}

func NewGetFollowingHandler(authMiddleware AuthMiddleware, followingRepo domain.FollowingRepository, router *httprouter.Router) *GetFollowingHandler {
	handler := &GetFollowingHandler{
		authMiddlerware: authMiddleware,
		followingRepo:   followingRepo,
		router:          router,
	}

	handler.router.GET("/api/v1/following", handler.authMiddlerware.Handle(handler.Handle()))
	return handler
}

func (handler GetFollowingHandler) Handle() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		user := (r.Context().Value(LoggedInUser)).(domain.User)

		result, err := handler.followingRepo.Following(r.Context(), user)

		if err != nil {
			JSONResponse(w, Error{err.Error()}, 500)
			return
		}

		JSONResponse(w, result, 200)
	}
}
