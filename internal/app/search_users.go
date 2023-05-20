package app

import (
	"mini-sns-ws/internal/domain"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type SearchUsersHandler struct {
	authMiddleware AuthMiddleware
	repo           domain.UserRepository
	router         *httprouter.Router
}

type searchUserResponse struct {
	Users []domain.User `json:"users"`
}

func NewSearchUsersHandler(
	authMiddleware AuthMiddleware,
	repo domain.UserRepository,
	router *httprouter.Router,
) *SearchUsersHandler {
	handler := &SearchUsersHandler{
		authMiddleware: authMiddleware,
		repo:           repo,
		router:         router,
	}

	handler.router.GET("/api/v1/search-users", CORS(handler.authMiddleware.Handle(handler.Handle())))

	return handler
}

func (handler *SearchUsersHandler) Handle() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		res := searchUserResponse{make([]domain.User, 0)}

		if r.URL.Query().Get("q") == "" {
			JSONResponse(w, res, 200)
			return
		}

		results, err := handler.repo.Search(r.Context(), r.URL.Query().Get("q"))

		if err != nil {
			JSONResponse(w, Error{err.Error()}, 422)
			return
		}

		JSONResponse(w, searchUserResponse{results}, 200)
	}
}
