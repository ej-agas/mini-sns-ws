package app

import (
	"mini-sns-ws/internal/domain"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const timeFormat = "January 2, 2006 3:04 PM MST"

type MyProfileHandler struct {
	authMiddleware AuthMiddleware
	followingRepo  domain.FollowingRepository
	postRepo       domain.PostRepository
	router         *httprouter.Router
}

func NewMyProfileHandler(
	authMiddleware AuthMiddleware,
	followingRepo domain.FollowingRepository,
	postRepo domain.PostRepository,
	router *httprouter.Router,
) *MyProfileHandler {
	handler := &MyProfileHandler{authMiddleware: authMiddleware, followingRepo: followingRepo, postRepo: postRepo, router: router}
	handler.router.GET("/api/v1/my-profile", CORS(handler.authMiddleware.Handle(handler.Handle())))

	return handler
}

type MyProfileResponse struct {
	ID             string `json:"id"`
	FullName       string `json:"full_name"`
	FirstName      string `json:"first_name"`
	MiddleName     string `json:"middle_name"`
	LastName       string `json:"last_name"`
	Bio            string `json:"bio"`
	Email          string `json:"email"`
	IsVerified     bool   `json:"is_verified"`
	VerifiedDate   string `json:"verified_date,omitempty"`
	JoinDate       string `json:"join_date"`
	FollowingCount int64  `json:"following_count"`
	FollowersCount int64  `json:"followers_count"`
	PostsCount     int64  `json:"posts_count"`
}

func (handler MyProfileHandler) Handle() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		user := (r.Context().Value(LoggedInUser)).(domain.User)

		var verifiedDate string
		if user.VerifiedAt.Time().Unix() != 0 {
			verifiedDate = user.VerifiedAt.Time().Format(timeFormat)
		}

		followingCount := handler.followingRepo.FollowingCount(r.Context(), user)
		followersCount := handler.followingRepo.FollowersCount(r.Context(), user)
		postsCount := handler.postRepo.PostsCount(r.Context(), user)

		response := MyProfileResponse{
			ID:             user.ID.Hex(),
			FullName:       user.FullName(),
			FirstName:      user.FirstName,
			MiddleName:     user.MiddleName,
			LastName:       user.LastName,
			Bio:            user.Bio,
			Email:          user.Email,
			IsVerified:     user.IsVerified,
			VerifiedDate:   verifiedDate,
			JoinDate:       user.CreatedAt.Time().Format(timeFormat),
			FollowingCount: followingCount,
			FollowersCount: followersCount,
			PostsCount:     postsCount,
		}

		JSONResponse(w, response, 200)
	}
}
