package usecase

import (
	"avito-winter-2025/internal/repo"
	myErrors "avito-winter-2025/internal/utils/errors"
	"context"
)

type CoinInterface interface {
	SendCoin(ctx context.Context, from uint32, to uint32, amount uint64) error
}

type Coin struct {
	repo repo.CoinInterface
}

func NewCoin(r repo.CoinInterface) CoinInterface {
	return &Coin{repo: r}
}

func (m *Coin) SendCoin(ctx context.Context, from uint32, to uint32, amount uint64) error {
	fromBalance, err := m.repo.CheckBalance(ctx, from)
	if err != nil {
		return err
	}
	if fromBalance < amount {
		return myErrors.NotEnoughCoinErr
	}
	err = m.repo.SendCoin(ctx, from, to, amount)
	if err != nil {
		return err
	}
	return nil
}
