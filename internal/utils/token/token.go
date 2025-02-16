package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint32 `json:"user_id"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}

type JWT struct {
	Secret  []byte
	ExpTime time.Duration
}

func NewJWT(secret string, duration string) (JWT, error) {
	expiration, err := time.ParseDuration(duration)
	if err != nil {
		return JWT{
			Secret:  []byte(secret),
			ExpTime: time.Duration(48 * time.Hour),
		}, err
	}
	return JWT{
		Secret:  []byte(secret),
		ExpTime: expiration,
	}, nil
}

func (j JWT) GenerateToken(userID uint32, name string) (string, error) {
	claims := Claims{
		UserID: userID,
		Name:   name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.ExpTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.Secret)
}
