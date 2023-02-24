package app

import (
	"mini-sns-ws/internal/domain"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type MyProfileHandler struct {
	authMiddleware AuthMiddleware
	router         *httprouter.Router
}

func NewMyProfileHandler(authMiddleware AuthMiddleware, router *httprouter.Router) *MyProfileHandler {
	handler := &MyProfileHandler{authMiddleware: authMiddleware, router: router}
	handler.router.GET("/api/v1/my-profile", handler.authMiddleware.Handle(handler.Handle()))

	return handler
}

func (handler MyProfileHandler) Handle() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		user := (r.Context().Value(LoggedInUser)).(domain.User)
		JSONResponse(w, user, 200)
	}
}
