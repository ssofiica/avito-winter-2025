package repo

import (
	"avito-winter-2025/internal/entity"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CoinInterface interface {
	SendCoin(ctx context.Context, transaction entity.Transaction) error
	CheckBalance(ctx context.Context, id uint32) (uint32, error)
	GetCoinHistory(ctx context.Context, id uint32) ([]entity.Transaction, error)
}

type Coin struct {
	db *pgxpool.Pool
}

func NewCoin(db *pgxpool.Pool) CoinInterface {
	return &Coin{db: db}
}

func (u *Coin) SendCoin(ctx context.Context, trans entity.Transaction) error {
	query := `insert into coin_history(from_user, to_user, amount, created_at) values ($1, $2, $3, NOW());`
	query1 := `update "user" set coins=coins-$1 where id=$2;`
	query2 := `update "user" set coins=coins+$1 where id=$2;`
	tx, err := u.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, query, trans.From, trans.To, trans.Amount)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, query1, trans.Amount, trans.From)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, query2, trans.Amount, trans.To)
	if err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func (u *Coin) CheckBalance(ctx context.Context, id uint32) (uint32, error) {
	query := `select coins from "user" where id=$1;`
	var res uint32
	err := u.db.QueryRow(ctx, query, id).Scan(&res)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (u *Coin) GetCoinHistory(ctx context.Context, id uint32) ([]entity.Transaction, error) {
	query := `select from_user, to_user, amount from coin_history where from_user=$1 OR to_user=$1;`
	var res []entity.Transaction
	rows, err := u.db.Query(ctx, query, id)
	if err != nil {
		return res, err
	}
	for rows.Next() {
		var t entity.Transaction
		err := rows.Scan(&t.From, &t.To, &t.Amount)
		if err != nil {
			return []entity.Transaction{}, err
		}
		res = append(res, t)
	}
	return res, nil
}
