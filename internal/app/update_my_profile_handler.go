package app

import (
	"encoding/json"
	"mini-sns-ws/internal/domain"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
)

type updateMyProfileInput struct {
	First_name  *string `json:"first_name" validate:"omitempty"`
	Middle_name *string `json:"middle_name" validate:"omitempty"`
	Last_name   *string `json:"last_name" validate:"omitempty"`
	Bio         *string `json:"bio" validate:"omitempty"`
	Password    *string `json:"password" validate:"omitempty,ascii,gte=8"`
}

type UpdateMyProfileHandler struct {
	authMiddleware AuthMiddleware
	hasher         Hasher
	validator      *validator.Validate
	userRepo       domain.UserRepository
	router         *httprouter.Router
}

func NewUpdateMyProfileHandler(authMiddleware AuthMiddleware, hasher Hasher, validator *validator.Validate, userRepo domain.UserRepository, router *httprouter.Router) *UpdateMyProfileHandler {
	handler := &UpdateMyProfileHandler{authMiddleware: authMiddleware, hasher: hasher, validator: validator, userRepo: userRepo, router: router}
	handler.router.PATCH("/api/v1/my-profile", CORS(handler.authMiddleware.Handle(handler.Handle())))

	return handler
}

type UpdateMyProfileResponse struct {
	FullName     string `json:"full_name"`
	Email        string `json:"email"`
	Bio          string `json:"bio"`
	IsVerified   bool   `json:"is_verified"`
	VerifiedDate string `json:"verified_date,omitempty"`
	JoinDate     string `json:"join_date"`
}

func (handler UpdateMyProfileHandler) Handle() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		input := updateMyProfileInput{}
		json.NewDecoder(r.Body).Decode(&input)
		errorResponse, err := Validate(handler.validator, input)

		if err != nil {
			JSONResponse(w, errorResponse, 422)
			return
		}
		user := (r.Context().Value(LoggedInUser)).(domain.User)

		if input.First_name != nil {
			user.FirstName = *input.First_name
		}

		if input.Middle_name != nil {
			user.MiddleName = *input.Middle_name
		}

		if input.Last_name != nil {
			user.LastName = *input.Last_name
		}

		if input.Password != nil {
			hashedPassword, err := handler.hasher.Hash(*input.Password)
			if err != nil {
				JSONResponse(w, Error{err.Error()}, 500)
				return
			}
			user.Password = hashedPassword
		}

		if input.Bio != nil {
			user.Bio = *input.Bio
		}

		if err := handler.userRepo.Save(r.Context(), user); err != nil {
			JSONResponse(w, Error{err.Error()}, 500)
			return
		}

		EmptyResponse(w, 200)
	}
}
