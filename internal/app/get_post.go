package app

import (
	"mini-sns-ws/internal/domain"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetPostHandler struct {
	repo   domain.PostRepository
	router *httprouter.Router
}

func NewGetPostHandler(repo domain.PostRepository, router *httprouter.Router) *GetPostHandler {
	handler := &GetPostHandler{repo: repo, router: router}

	handler.routes()

	return handler
}

func (handler *GetPostHandler) routes() {
	handler.router.GET("/posts/:id", handler.handle())
}

func (handler *GetPostHandler) handle() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		postObjId, err := primitive.ObjectIDFromHex(ps.ByName("id"))

		if err != nil {
			EmptyResponse(w, 404)
			return
		}

		post, err := handler.repo.FindOneBy(r.Context(), "_id", postObjId)

		if err != nil {
			EmptyResponse(w, 404)
		}

		JSONResponse(w, post, 200)
	}
}
