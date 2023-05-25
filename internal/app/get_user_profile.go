package app

import (
	"mini-sns-ws/internal/domain"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const humanDateFormat = "January 2, 2006 3:04 PM MST"

type UserProfileHandler struct {
	authMiddleware AuthMiddleware
	userRepo       domain.UserRepository
	router         *httprouter.Router
}

func NewUserProfileHandler(authMiddleware AuthMiddleware, userRepo domain.UserRepository, router *httprouter.Router) *UserProfileHandler {
	handler := &UserProfileHandler{authMiddleware: authMiddleware, userRepo: userRepo, router: router}
	handler.router.GET("/api/v1/users/:id", CORS(handler.authMiddleware.Handle(handler.Handle())))

	return handler
}

type UserProfileResponse struct {
	FullName     string `json:"full_name"`
	FirstName    string `json:"first_name"`
	MiddleName   string `json:"middle_name"`
	LastName     string `json:"last_name"`
	Bio          string `json:"bio"`
	Email        string `json:"email"`
	IsVerified   bool   `json:"is_verified"`
	VerifiedDate string `json:"verified_date,omitempty"`
	JoinDate     string `json:"join_date"`
}

func (handler UserProfileHandler) Handle() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		userId, err := primitive.ObjectIDFromHex(ps.ByName("id"))
		if err != nil {
			JSONResponse(w, Error{ErrInvalidUserId.Error()}, 422)
			return
		}

		filter := domain.NewFilter()
		filter["_id"] = userId

		user, err := handler.userRepo.FindOneBy(r.Context(), filter, *domain.NewFindOptions())

		if err != nil {
			JSONResponse(w, Error{err.Error()}, 422)
			return
		}

		var verifiedDate string
		if user.VerifiedAt.Time().Unix() != 0 {
			verifiedDate = user.VerifiedAt.Time().Format(humanDateFormat)
		}

		response := MyProfileResponse{
			FullName:     user.FullName(),
			FirstName:    user.FirstName,
			MiddleName:   user.MiddleName,
			LastName:     user.LastName,
			Bio:          user.Bio,
			Email:        user.Email,
			IsVerified:   user.IsVerified,
			VerifiedDate: verifiedDate,
			JoinDate:     user.CreatedAt.Time().Format(humanDateFormat),
		}

		JSONResponse(w, response, 200)
	}
}
