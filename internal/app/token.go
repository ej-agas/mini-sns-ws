package app

import (
	"mini-sns-ws/internal/domain"
	"time"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
)

type TokenService interface {
	GenerateFor(user domain.User) (string, error)
	IsExpired(token jwt.Token) bool
	FromString(token string) (jwt.Token, error)
}

type JWTTokenService struct {
	SecretKey string
	Expiry    time.Duration
}

func (service JWTTokenService) GenerateFor(user domain.User) (string, error) {
	now := time.Now()
	token, err := jwt.NewBuilder().
		Issuer("mini-sns-ws").
		Claim("user_id", user.ID.Hex()).
		IssuedAt(now).
		Expiration(time.Now().Add(service.Expiry)).
		Build()

	if err != nil {
		return "", err
	}

	serialized, err := jwt.Sign(token, jwa.HS512, []byte(service.SecretKey))

	if err != nil {
		return "", err
	}

	return string(serialized), nil
}

func (service JWTTokenService) IsExpired(token jwt.Token) bool {
	return time.Now().After(token.Expiration())
}

func (service JWTTokenService) FromString(tokenAsString string) (jwt.Token, error) {
	token, err := jwt.Parse([]byte(tokenAsString), jwt.WithVerify(jwa.HS512, []byte(service.SecretKey)), jwt.WithValidate(true))

	if err != nil {
		return nil, err
	}

	return token, nil
}
