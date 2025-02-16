package delivery

import (
	"avito-winter-2025/internal/entity"
	"avito-winter-2025/internal/usecase"
	myErrors "avito-winter-2025/internal/utils/errors"
	"avito-winter-2025/internal/utils/response"
	"errors"

	"context"
	"net/http"

	"github.com/gorilla/mux"
)

type ShopHandler struct {
	merchUC usecase.MerchInterface
	userUC  usecase.UserInterface
	coinUC  usecase.CoinInterface
}

func NewShopHandler(m usecase.MerchInterface, u usecase.UserInterface, c usecase.CoinInterface) *ShopHandler {
	return &ShopHandler{merchUC: m, userUC: u, coinUC: c}
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
		response.WithError(w, 400, ErrNoRequestVars)
		return
	}
	err := h.merchUC.Buy(context.Background(), user.ID, merchName)
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

func (h *ShopHandler) GetInfo(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(userKey).(entity.User)
	if !ok {
		response.WithError(w, 401, ErrDefault401)
		return
	}
	u, err := h.userUC.GetUser(context.Background(), "", user.ID)
	if err != nil {
		if errors.Is(err, myErrors.NoUserErr) {
			response.WithError(w, 400, myErrors.NoUserErr)
			return
		}
		response.WithError(w, 500, ErrDefault500)
		return
	}
	inventory, err := h.merchUC.GetInventoryHistory(context.Background(), user.ID)
	if err != nil {
		response.WithError(w, 500, ErrDefault500)
		return
	}
	coinHistory, err := h.coinUC.GetCoinHistory(context.Background(), user.ID)
	if err != nil {
		response.WithError(w, 500, ErrDefault500)
		return
	}
	res := entity.InfoResponse{
		Coins:       u.Coins,
		Inventory:   inventory,
		CoinHistory: coinHistory,
	}
	response.WriteData(w, res, 200)
}
