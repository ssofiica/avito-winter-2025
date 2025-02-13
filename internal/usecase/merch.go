package usecase

import (
	"avito-winter-2025/internal/repo"
	myErrors "avito-winter-2025/internal/utils/errors"
	"context"
)

type MerchInterface interface {
	Buy(ctx context.Context, userId uint32, merchName string) error
}

type Merch struct {
	merchRepo repo.MerchInterface
	coinRepo  repo.CoinInterface
}

func NewMerch(m repo.MerchInterface, c repo.CoinInterface) MerchInterface {
	return &Merch{merchRepo: m, coinRepo: c}
}

func (m *Merch) Buy(ctx context.Context, userId uint32, merchName string) error {
	merch, err := m.merchRepo.GetByName(ctx, merchName)
	if err != nil {
		return err
	}
	if merch == nil {
		return myErrors.NoMerchErr
	}
	balance, err := m.coinRepo.CheckBalance(ctx, userId)
	if err != nil {
		return err
	}
	if balance < merch.Cost {
		return myErrors.NotEnoughCoinErr
	}
	err = m.merchRepo.Buy(ctx, userId, merch.ID, merch.Cost)
	if err != nil {
		return err
	}
	return nil
}
