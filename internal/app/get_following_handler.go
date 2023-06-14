package app

import (
	"mini-sns-ws/internal/domain"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Following struct {
	ID      primitive.ObjectID `json:"_id"`
	Name    string             `json:"name"`
	Bio     string             `json:"bio"`
	Picture string             `json:"picture"`
}

type GetFollowingHandler struct {
	authMiddlerware AuthMiddleware
	followingRepo   domain.FollowingRepository
	userRepo        domain.UserRepository
	router          *httprouter.Router
}

func NewGetFollowingHandler(authMiddleware AuthMiddleware, followingRepo domain.FollowingRepository, userRepo domain.UserRepository, router *httprouter.Router) *GetFollowingHandler {
	handler := &GetFollowingHandler{
		authMiddlerware: authMiddleware,
		followingRepo:   followingRepo,
		userRepo:        userRepo,
		router:          router,
	}

	handler.router.GET("/api/v1/following", handler.authMiddlerware.Handle(handler.Handle()))
	return handler
}

func (handler GetFollowingHandler) Handle() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		user := (r.Context().Value(LoggedInUser)).(domain.User)

		followingIds := make([]primitive.ObjectID, 0)
		followingUsers := make([]Following, 0)

		followingResult, err := handler.followingRepo.Following(r.Context(), user)

		if err != nil {
			JSONResponse(w, Error{err.Error()}, 500)
			return
		}

		for _, following := range followingResult {
			followingIds = append(followingIds, following.Following)
		}

		filter := domain.NewFilter()
		filter["_id"] = bson.M{"$in": followingIds}

		users, err := handler.userRepo.FindBy(r.Context(), filter, *domain.NewFindOptions())

		if err != nil {
			JSONResponse(w, Error{err.Error()}, 500)
			return
		}

		for _, user := range users {
			following := Following{
				ID:      user.ID,
				Name:    user.FullName(),
				Bio:     user.Bio,
				Picture: user.Picture,
			}
			followingUsers = append(followingUsers, following)
		}

		JSONResponse(w, struct {
			Data []Following `json:"data"`
		}{Data: followingUsers},
			200)
	}
}
