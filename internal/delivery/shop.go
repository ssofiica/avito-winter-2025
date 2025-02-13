package delivery

import (
	"avito-winter-2025/internal/entity"
	"avito-winter-2025/internal/usecase"
	myErrors "avito-winter-2025/internal/utils/errors"
	"avito-winter-2025/internal/utils/response"
	"errors"
	"fmt"

	"context"
	"net/http"

	"github.com/gorilla/mux"
)

type ShopHandler struct {
	usecase usecase.MerchInterface
}

func NewShopHandler(u usecase.MerchInterface) *ShopHandler {
	return &ShopHandler{usecase: u}
}

func (h *ShopHandler) BuyMerch(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(userKey).(entity.User)
	if !ok {
		response.WithError(w, 401, ErrDefault401)
		return
	}
	vars := mux.Vars(r)
	merchName := vars["item"]
	if merchName == "" {
		fmt.Println("no item")
		response.WithError(w, 400, ErrNoRequestVars)
		return
	}
	err := h.usecase.Buy(context.Background(), user.ID, merchName)
	if err != nil {
		if errors.Is(err, myErrors.NotEnoughCoinErr) {
			response.WithError(w, 400, myErrors.NotEnoughCoinErr)
			return
		}
		if errors.Is(err, myErrors.NoMerchErr) {
			response.WithError(w, 400, myErrors.NoMerchErr)
			return
		}
		response.WithError(w, 500, ErrDefault500)
		return
	}
	response.WriteData(w, nil, 200)
}
