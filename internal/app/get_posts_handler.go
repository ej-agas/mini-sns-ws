package app

import (
	"mini-sns-ws/internal/domain"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type GetPostsHandler struct {
	authMiddlerware AuthMiddleware
	repo            domain.PostRepository
	router          *httprouter.Router
}

func NewGetPostsHandler(authMiddleware AuthMiddleware, repo domain.PostRepository, router *httprouter.Router) *GetPostsHandler {
	handler := &GetPostsHandler{authMiddlerware: authMiddleware, repo: repo, router: router}
	handler.router.GET("/api/v1/posts", handler.authMiddlerware.Handle(handler.GetPosts()))

	return handler
}

func (handler GetPostsHandler) GetPosts() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		user := r.Context().Value(LoggedInUser).(domain.User)

		filter := domain.NewFilter()
		filter["user_id"] = user.ID
		result, err := handler.repo.FindBy(r.Context(), filter)

		if err != nil {
			EmptyResponse(w, 404)
			return
		}

		JSONResponse(w, domain.NewModelCollection(result), 200)
	}
}
