package app

import (
	"mini-sns-ws/internal/domain"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetPostHandler struct {
	authMiddleware AuthMiddleware
	repo           domain.PostRepository
	router         *httprouter.Router
}

func NewGetPostHandler(authMiddleware AuthMiddleware, repo domain.PostRepository, router *httprouter.Router) *GetPostHandler {
	handler := &GetPostHandler{authMiddleware: authMiddleware, repo: repo, router: router}
	handler.router.GET("/api/v1/posts/:id", CORS(handler.authMiddleware.Handle(handler.handle())))

	return handler
}

func (handler *GetPostHandler) handle() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		postObjId, err := primitive.ObjectIDFromHex(ps.ByName("id"))

		if err != nil {
			EmptyResponse(w, 404)
			return
		}

		filter := domain.NewFilter()
		filter["_id"] = postObjId

		post, err := handler.repo.FindOneBy(r.Context(), filter)

		if err != nil {
			EmptyResponse(w, 404)
			return
		}

		JSONResponse(w, post, 200)
	}
}
