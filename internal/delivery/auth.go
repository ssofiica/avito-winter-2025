package delivery

import (
	"avito-winter-2025/internal/entity"
	"avito-winter-2025/internal/usecase"
	myErrors "avito-winter-2025/internal/utils/errors"
	"avito-winter-2025/internal/utils/request"
	"avito-winter-2025/internal/utils/response"
	"avito-winter-2025/internal/utils/token"
	"errors"
	"fmt"

	"context"
	"net/http"
)

type AuthHandler struct {
	usecase usecase.UserInterface
}

func NewAuthHandler(u usecase.UserInterface) *AuthHandler {
	return &AuthHandler{usecase: u}
}

func (h *AuthHandler) Auth(w http.ResponseWriter, r *http.Request) {
	payload := entity.AuthRequest{}
	if err := request.GetRequestData(r, &payload); err != nil {
		response.WithError(w, 400, ErrDefault400)
		return
	}
	if payload.Name == "" || payload.Password == "" {
		response.WithError(w, 400, ErrDefault400)
		return
	}
	userData, err := h.usecase.Auth(context.Background(), payload)
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, myErrors.WrongLoginOrPasswordErr) {
			response.WithError(w, 500, myErrors.WrongLoginOrPasswordErr)
			return
		}
		response.WithError(w, 500, ErrDefault500)
		return
	}
	jwtToken, err := token.GenerateToken(userData.ID, userData.Name)
	if err != nil {
		response.WithError(w, 500, ErrTokenGenerate)
		return
	}
	w.Header().Set("Authorization", "Bearer "+jwtToken)
	res := entity.AuthResponse{Token: jwtToken}
	response.WriteData(w, res, 200)
}
