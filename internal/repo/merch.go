package repo

import (
	"avito-winter-2025/internal/entity"
	// myErrors "avito-winter-2025/internal/utils/errors"
	"context"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MerchInterface interface {
	Buy(ctx context.Context, userId uint32, merchId uint32, cost uint32) error
	GetByName(ctx context.Context, name string) (*entity.Merch, error)
}

type Merch struct {
	db *pgxpool.Pool
}

func NewMerch(db *pgxpool.Pool) MerchInterface {
	return &Merch{db: db}
}

func (u *Merch) Buy(ctx context.Context, userId uint32, merchId uint32, cost uint32) error {
	queryInsert := `insert into inventory(merch_id, user_id, created_at) values ($1, $2, NOW())`
	queryUpdate := `update "user" set coins=coins-$1 where id=$2;`
	tx, err := u.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, queryInsert, merchId, userId)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, queryUpdate, cost, userId)
	if err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func (u *Merch) GetByName(ctx context.Context, name string) (*entity.Merch, error) {
	query := `select id, name, cost from merch where name=$1`
	var res entity.Merch
	err := u.db.QueryRow(ctx, query, name).Scan(&res.ID, &res.Name, &res.Cost)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return nil, nil
		}
		return nil, err
	}
	return &res, nil
}
