package app

import (
	"mini-sns-ws/internal/domain"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DeletePostHandler struct {
	authMiddleware AuthMiddleware
	repository     domain.PostRepository
	router         *httprouter.Router
}

func NewDeletePostHandler(authMiddleware AuthMiddleware, repository domain.PostRepository, router *httprouter.Router) *DeletePostHandler {
	handler := &DeletePostHandler{authMiddleware: authMiddleware, repository: repository, router: router}

	handler.router.DELETE("/api/v1/posts/:id", handler.authMiddleware.Handle(handler.DeletePost()))

	return handler
}

func (handler DeletePostHandler) DeletePost() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		user := r.Context().Value(LoggedInUser).(domain.User)

		postId, err := primitive.ObjectIDFromHex(ps.ByName("id"))

		if err != nil {
			JSONResponse(w, invalidPostId, 422)
			return
		}

		filter := domain.NewFilter()
		filter["_id"] = postId
		filter["user_id"] = user.ID

		if err := handler.repository.DeleteBy(r.Context(), filter); err != nil {
			JSONResponse(w, Error{Message: err.Error()}, 400)
			return
		}

		EmptyResponse(w, 200)
	}
}
