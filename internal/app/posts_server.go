package app

import (
	"encoding/json"
	"mini-sns-ws/internal/domain"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type PostsServer struct {
	repo   domain.UserRepository
	router *httprouter.Router
}

func NewPostsServer(userRepo domain.UserRepository, r *httprouter.Router) *PostsServer {
	ps := &PostsServer{repo: userRepo, router: r}
	ps.routes()

	return ps
}

func (ps *PostsServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ps.router.ServeHTTP(w, r)
}

func (ps *PostsServer) handlePosts() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		res := struct {
			Message string
		}{Message: "Posts Service"}

		json.NewEncoder(w).Encode(res)
	}
}
