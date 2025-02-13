package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

//var COINS = 1000

type CoinInterface interface {
	SendCoin(ctx context.Context, from uint32, to uint32, amount uint64) error
	CheckBalance(ctx context.Context, id uint32) (uint64, error)
}

type Coin struct {
	db *pgxpool.Pool
}

func NewCoin(db *pgxpool.Pool) CoinInterface {
	return &Coin{db: db}
}

func (u *Coin) SendCoin(ctx context.Context, from uint32, to uint32, amount uint64) error {
	query := `insert into coin_history(from_user, to_user, amount) values ($1, $2, $3);`
	query1 := `update "user" set coins=coins-$1 where id=$2;`
	query2 := `update "user" set coins=coins+$1 where id=$2;`
	tx, err := u.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, query, from, to, amount)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, query1, amount, from)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, query2, amount, to)
	if err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func (u *Coin) CheckBalance(ctx context.Context, id uint32) (uint64, error) {
	query := `select coins from "user" where id=$1;`
	var res uint64
	err := u.db.QueryRow(ctx, query, id).Scan(&res)
	if err != nil {
		return 0, err
	}
	return res, nil
}
