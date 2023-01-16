package app

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type contextKey string

const LoggedInUser contextKey = "user"

func AuthMiddleware(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if bearerToken := r.Header.Get("Authorization"); bearerToken == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		// ctx = context.WithValue(ctx, LoggedInUser, domain.User{FirstName: "Jane", LastName: "doe"})

		next(w, r.WithContext(ctx), ps)
	}
}
