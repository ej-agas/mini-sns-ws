package app

import (
	"fmt"
	"mini-sns-ws/internal/domain"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const humanDateFormat = "January 2, 2006 3:04 PM MST"

type UserProfileHandler struct {
	authMiddleware AuthMiddleware
	userRepo       domain.UserRepository
	postRepo       domain.PostRepository
	followingRepo  domain.FollowingRepository
	router         *httprouter.Router
}

func NewUserProfileHandler(authMiddleware AuthMiddleware, userRepo domain.UserRepository, postRepo domain.PostRepository, followingRepo domain.FollowingRepository, router *httprouter.Router) *UserProfileHandler {
	handler := &UserProfileHandler{authMiddleware: authMiddleware, userRepo: userRepo, postRepo: postRepo, followingRepo: followingRepo, router: router}
	handler.router.GET("/api/v1/users/:id", CORS(handler.authMiddleware.Handle(handler.Handle())))

	return handler
}

type UserProfileResponse struct {
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
	IsFollowing    bool   `json:"is_following"`
	FollowingCount int64  `json:"following_count"`
	FollowersCount int64  `json:"followers_count"`
	PostsCount     int64  `json:"posts_count"`
}

func (handler UserProfileHandler) Handle() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		loggedInUser := (r.Context().Value(LoggedInUser)).(domain.User)
		filter := domain.NewFilter()

		userId, err := primitive.ObjectIDFromHex(ps.ByName("id"))

		if err != nil {
			filter["username"] = ps.ByName("id")
		} else {
			filter["_id"] = userId
		}

		user, err := handler.userRepo.FindOneBy(r.Context(), filter, *domain.NewFindOptions())

		if err != nil {
			JSONResponse(w, Error{"user not found."}, 404)
			return
		}

		fmt.Println(123)

		postsCount := handler.postRepo.PostsCount(r.Context(), user)
		followersCount := handler.followingRepo.FollowersCount(r.Context(), user)
		followingCount := handler.followingRepo.FollowingCount(r.Context(), user)
		isFollowing := handler.followingRepo.IsFollowing(r.Context(), loggedInUser, user)
		var verifiedDate string
		if user.VerifiedAt.Time().Unix() != 0 {
			verifiedDate = user.VerifiedAt.Time().Format(humanDateFormat)
		}

		response := UserProfileResponse{
			ID:             user.ID.Hex(),
			FullName:       user.FullName(),
			FirstName:      user.FirstName,
			MiddleName:     user.MiddleName,
			LastName:       user.LastName,
			Bio:            user.Bio,
			Email:          user.Email,
			IsVerified:     user.IsVerified,
			VerifiedDate:   verifiedDate,
			JoinDate:       user.CreatedAt.Time().Format(humanDateFormat),
			IsFollowing:    isFollowing,
			FollowingCount: followingCount,
			FollowersCount: followersCount,
			PostsCount:     postsCount,
		}

		JSONResponse(w, response, 200)
	}
}
