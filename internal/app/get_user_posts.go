package app

import (
	"mini-sns-ws/internal/domain"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetUserPostsHandler struct {
	authMiddlerware AuthMiddleware
	repo            domain.PostRepository
	router          *httprouter.Router
}

func NewGetUserPostsHandler(authMiddleware AuthMiddleware, repo domain.PostRepository, router *httprouter.Router) *GetUserPostsHandler {
	handler := &GetUserPostsHandler{authMiddlerware: authMiddleware, repo: repo, router: router}
	handler.router.GET("/api/v1/users/:id/posts", CORS(handler.authMiddlerware.Handle(handler.Handle())))

	return handler
}

func (handler GetUserPostsHandler) Handle() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		filter := make(map[string]interface{}, 1)

		userId, err := primitive.ObjectIDFromHex(ps.ByName("id"))
		if err != nil {
			filter["username"] = ps.ByName("id")
			result, err := handler.repo.FindBy(r.Context(), filter, *domain.NewFindOptions())

			if err != nil {
				JSONResponse(w, Error{err.Error()}, 422)
				return
			}

			JSONResponse(w, domain.NewModelCollection(result), 200)
			return
		}

		filter["user_id"] = userId

		result, err := handler.repo.FindBy(r.Context(), filter, *domain.NewFindOptions())

		if err != nil {
			EmptyResponse(w, 404)
			return
		}

		JSONResponse(w, domain.NewModelCollection(result), 200)
	}
}
