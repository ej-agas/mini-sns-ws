package app

import (
	"encoding/json"
	"mini-sns-ws/internal/domain"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type UserServer struct {
	repo   domain.UserRepository
	router *httprouter.Router
}

func NewUserServer(userRepo domain.UserRepository, r *httprouter.Router) *UserServer {
	s := &UserServer{repo: userRepo, router: r}
	s.routes()

	return s
}

func (s *UserServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *UserServer) handleRegister() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Header().Add("content-type", "application/json; charset=utf-8")

		res := struct {
			Message string
		}{Message: "User Service"}

		json.NewEncoder(w).Encode(res)
	}
}
