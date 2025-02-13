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

var JwtSecret = []byte("your_secret_key")

func GenerateToken(userID uint32, name string) (string, error) {
	claims := Claims{
		UserID: userID,
		Name:   name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(48 * time.Hour)), // Токен истекает через 48 часов
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtSecret)
}
