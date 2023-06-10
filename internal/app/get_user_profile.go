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
	followingRepo  domain.FollowingRepository
	router         *httprouter.Router
}

func NewUserProfileHandler(authMiddleware AuthMiddleware, userRepo domain.UserRepository, router *httprouter.Router) *UserProfileHandler {
	handler := &UserProfileHandler{authMiddleware: authMiddleware, userRepo: userRepo, router: router}
	handler.router.GET("/api/v1/users/:id", CORS(handler.authMiddleware.Handle(handler.Handle())))

	return handler
}

type UserProfileResponse struct {
	ID           string `json:"id"`
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

		loggedInUser := (r.Context().Value(LoggedInUser)).(domain.User)

		userId, err := primitive.ObjectIDFromHex(ps.ByName("id"))
		if err != nil {

			filter := domain.NewFilter()
			filter["username"] = ps.ByName("id")

			user, err := handler.userRepo.FindOneBy(r.Context(), filter, *domain.NewFindOptions())

			if err != nil {
				JSONResponse(w, Error{err.Error()}, 422)
				return
			}

			JSONResponse(w, createProfileResponse(user), 200)

			return
		}

		filter := domain.NewFilter()
		filter["_id"] = userId

		user, err := handler.userRepo.FindOneBy(r.Context(), filter, *domain.NewFindOptions())

		if err != nil {
			JSONResponse(w, Error{err.Error()}, 422)
			return
		}

		JSONResponse(w, createProfileResponse(user), 200)
	}
}

func createProfileResponse(u domain.User) UserProfileResponse {
	var verifiedDate string
	if u.VerifiedAt.Time().Unix() != 0 {
		verifiedDate = u.VerifiedAt.Time().Format(humanDateFormat)
	}

	return UserProfileResponse{
		ID:           u.ID.Hex(),
		FullName:     u.FullName(),
		FirstName:    u.FirstName,
		MiddleName:   u.MiddleName,
		LastName:     u.LastName,
		Bio:          u.Bio,
		Email:        u.Email,
		IsVerified:   u.IsVerified,
		VerifiedDate: verifiedDate,
		JoinDate:     u.CreatedAt.Time().Format(humanDateFormat),
	}
}
