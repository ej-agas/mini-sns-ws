package app

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (h *UserHandler) verify() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		token := r.URL.Query().Get("token")

		if token == "" {
			EmptyResponse(w, 404)
			return
		}

		userId, err := h.keyValueStore.Get(r.Context(), token)

		if err != nil {
			JSONResponse(w, Error{"Invalid or expired token."}, 422)
			return
		}

		user, err := h.repo.Find(r.Context(), userId)

		if err != nil {
			JSONResponse(w, err.Error(), 400)
			return
		}

		user.Verify()
		if err := h.repo.Save(r.Context(), user); err != nil {
			JSONResponse(w, Error{fmt.Sprintf("Error verifying user: %s", err.Error())}, 400)
			return
		}

		h.keyValueStore.Delete(r.Context(), token)

		EmptyResponse(w, 200)
	}
}
