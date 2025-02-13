package usecase

import (
	"avito-winter-2025/internal/entity"
	"avito-winter-2025/internal/repo"
	myErrors "avito-winter-2025/internal/utils/errors"
	"context"
)

type CoinInterface interface {
	SendCoin(ctx context.Context, from uint32, to uint32, amount uint32) error
	GetCoinHistory(ctx context.Context, id uint32) (entity.CoinHistory, error)
}

type Coin struct {
	coinRepo repo.CoinInterface
	userRepo repo.UserInterface
}

func NewCoin(c repo.CoinInterface, u repo.UserInterface) CoinInterface {
	return &Coin{coinRepo: c, userRepo: u}
}

func (u *Coin) SendCoin(ctx context.Context, from uint32, to uint32, amount uint32) error {
	fromBalance, err := u.coinRepo.CheckBalance(ctx, from)
	if err != nil {
		return err
	}
	if fromBalance < amount {
		return myErrors.NotEnoughCoinErr
	}
	err = u.coinRepo.SendCoin(ctx, entity.Transaction{From: from, To: to, Amount: amount})
	if err != nil {
		return err
	}
	return nil
}

func (u *Coin) GetCoinHistory(ctx context.Context, id uint32) (entity.CoinHistory, error) {
	received := []entity.Received{}
	sent := []entity.Sent{}
	res, err := u.coinRepo.GetCoinHistory(ctx, id)
	if err != nil {
		return entity.CoinHistory{Received: received, Sent: sent}, err
	}
	for _, trans := range res {
		if trans.From == id {
			toUser, err := u.userRepo.GetUser(ctx, "", trans.To)
			if err != nil {
				return entity.CoinHistory{Received: received, Sent: sent}, err
			}
			sent = append(sent, entity.Sent{
				ToUser: toUser.Name,
				Amount: trans.Amount,
			})
			continue
		}
		if trans.To == id {
			fromUser, err := u.userRepo.GetUser(ctx, "", trans.From)
			if err != nil {
				return entity.CoinHistory{Received: received, Sent: sent}, err
			}
			received = append(received, entity.Received{
				FromUser: fromUser.Name,
				Amount:   trans.Amount,
			})
		}
	}
	return entity.CoinHistory{Received: received, Sent: sent}, nil
}
