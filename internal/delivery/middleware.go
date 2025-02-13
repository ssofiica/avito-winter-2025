package delivery

import (
	"avito-winter-2025/internal/entity"
	"avito-winter-2025/internal/utils/response"
	"avito-winter-2025/internal/utils/token"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

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
			// Проверка метода подписи
			if _, ok := tok.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("недопустимый метод подписи")
			}
			return token.JwtSecret, nil
		})
		if err != nil {
			response.WithError(w, 401, ErrDefault401)
			return
		}
		if !t.Valid {
			fmt.Println("Недействительный токен")
			response.WithError(w, 401, ErrDefault401)
		}

		// Проверяем срок действия токена
		if claims.ExpiresAt.Time.Before(time.Now()) {
			response.WithError(w, 401, ErrDefault401)
			return
		}
		// Добавляем данные пользователя в контекст (можно использовать context для передачи дальше)
		user := entity.User{ID: claims.UserID, Name: claims.Name}
		var key string = "user"
		ctx := context.WithValue(r.Context(), key, user)
		r = r.WithContext(ctx)
		// Передаем управление следующему обработчику
		next.ServeHTTP(w, r)
	}
}
