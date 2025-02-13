package delivery

import (
	"avito-winter-2025/internal/entity"
	"avito-winter-2025/internal/usecase"
	myErrors "avito-winter-2025/internal/utils/errors"
	"avito-winter-2025/internal/utils/request"
	"avito-winter-2025/internal/utils/response"
	"errors"
	"fmt"

	"context"
	"net/http"
)

type CoinHandler struct {
	coinUC usecase.CoinInterface
	userUC usecase.UserInterface
}

func NewCoinHandler(c usecase.CoinInterface, u usecase.UserInterface) *CoinHandler {
	return &CoinHandler{coinUC: c, userUC: u}
}

func (h *CoinHandler) SendCoin(w http.ResponseWriter, r *http.Request) {
	from, ok := r.Context().Value("user").(entity.User)
	if !ok {
		response.WithError(w, 401, ErrDefault401)
		return
	}
	payload := entity.SendCoinRequest{}
	if err := request.GetRequestData(r, &payload); err != nil {
		response.WithError(w, 400, ErrDefault400)
		return
	}
	if !payload.Valid() {
		response.WithError(w, 400, ErrDefault400)
		return
	}
	to, err := h.userUC.GetUser(context.Background(), payload.ToUser, 0)
	if err != nil {
		if errors.Is(err, myErrors.NoUserErr) {
			response.WithError(w, 400, myErrors.NoUserErr)
			return
		}
		response.WithError(w, 500, ErrDefault500)
		return
	}
	err = h.coinUC.SendCoin(context.Background(), from.ID, to.ID, uint64(payload.Amount))
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, myErrors.NotEnoughCoinErr) {
			response.WithError(w, 400, myErrors.NotEnoughCoinErr)
			return
		}
		response.WithError(w, 500, ErrDefault500)
		return
	}
	response.WriteData(w, nil, 200)
}
