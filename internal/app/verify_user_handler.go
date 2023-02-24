package app

import (
	"fmt"
	"mini-sns-ws/internal/domain"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type VerifyUserHandler struct {
	repo          domain.UserRepository
	tokenService  TokenService
	keyValueStore domain.KeyValueStore
	router        *httprouter.Router
}

func NewVerifyUserHandler(
	userRepo domain.UserRepository,
	tokenService TokenService,
	keyValueStore domain.KeyValueStore,
	router *httprouter.Router,
) *VerifyUserHandler {
	handler := &VerifyUserHandler{
		repo:          userRepo,
		tokenService:  tokenService,
		keyValueStore: keyValueStore,
		router:        router,
	}

	handler.router.GET("/api/v1/verify", handler.verify())

	return handler
}

type verifyUserResponse struct {
	Message string `json:"token"`
}

func (handler *VerifyUserHandler) verify() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		token := r.URL.Query().Get("token")

		if token == "" {
			JSONResponse(w, Error{"Invalid or expired token."}, 422)
			return
		}

		userId, err := handler.keyValueStore.Get(r.Context(), token)

		if err != nil {
			JSONResponse(w, Error{"Invalid or expired token."}, 422)
			return
		}

		user, err := handler.repo.Find(r.Context(), userId)

		if err != nil {
			JSONResponse(w, err.Error(), 400)
			return
		}

		user.Verify()

		if err := handler.repo.Save(r.Context(), user); err != nil {
			JSONResponse(w, Error{fmt.Sprintf("Error verifying user: %s", err.Error())}, 400)
			return
		}

		handler.keyValueStore.Delete(r.Context(), token)
		jwtToken, err := handler.tokenService.GenerateFor(user)

		if err != nil {
			JSONResponse(w, Error{err.Error()}, 400)
			return
		}

		JSONResponse(w, verifyUserResponse{jwtToken}, 200)
	}
}
