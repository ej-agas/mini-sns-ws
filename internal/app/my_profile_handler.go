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

type MyProfileResponse struct {
	FullName     string `json:"full_name"`
	Email        string `json:"email"`
	IsVerified   bool   `json:"is_verified"`
	VerifiedDate string `json:"verified_date,omitempty"`
	JoinDate     string `json:"join_date"`
}

func (handler MyProfileHandler) Handle() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		user := (r.Context().Value(LoggedInUser)).(domain.User)

		var verifiedDate string
		if user.VerifiedAt.Time().Unix() != 0 {
			verifiedDate = user.VerifiedAt.Time().Format("January 2, 2006 3:04 PM MST")
		}

		response := MyProfileResponse{
			FullName:     user.FullName(),
			Email:        user.Email,
			IsVerified:   user.IsVerified,
			VerifiedDate: verifiedDate,
			JoinDate:     user.CreatedAt.Time().Format("January 2, 2006 3:04 PM MST"),
		}

		JSONResponse(w, response, 200)
	}
}
