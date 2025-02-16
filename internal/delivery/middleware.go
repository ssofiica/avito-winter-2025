package delivery

import (
	"avito-winter-2025/internal/entity"
	"avito-winter-2025/internal/utils/response"
	"avito-winter-2025/internal/utils/token"
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var userKey string = "user"

func JWTMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			fmt.Println(ErrDefault401.Error())
			response.WithError(w, 401, ErrDefault401)
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			fmt.Println("Неверный формат токена")
			response.WithError(w, 401, ErrDefault401)
			return
		}
		tokenString := parts[1]

		claims := &token.Claims{}
		t, err := jwt.ParseWithClaims(tokenString, claims, func(tok *jwt.Token) (interface{}, error) {
			if _, ok := tok.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("недопустимый метод подписи")
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil {
			response.WithError(w, 401, ErrDefault401)
			return
		}
		if !t.Valid {
			fmt.Println("Недействительный токен")
			response.WithError(w, 401, ErrDefault401)
		}

		if claims.ExpiresAt.Time.Before(time.Now()) {
			response.WithError(w, 401, ErrDefault401)
			return
		}

		user := entity.User{ID: claims.UserID, Name: claims.Name}
		ctx := context.WithValue(r.Context(), userKey, user)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}
}
