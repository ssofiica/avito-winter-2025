package repo

import (
	"avito-winter-2025/internal/entity"
	"context"

	"github.com/jackc/pgx"
)

//go:generate mockgen -source=merch.go -destination=mock/merch_mock.go -package=mock
type MerchInterface interface {
	Buy(ctx context.Context, userId uint32, merchId uint32, cost uint32) error
	GetByName(ctx context.Context, name string) (*entity.Merch, error)
	GetInventoryHistory(ctx context.Context, id uint32) ([]entity.Inventory, error)
}

type Merch struct {
	db DBInterface
}

func NewMerch(db DBInterface) MerchInterface {
	return &Merch{db: db}
}

func (m *Merch) Buy(ctx context.Context, userId uint32, merchId uint32, cost uint32) error {
	queryInsert := `insert into inventory(merch_id, user_id, created_at) values ($1, $2, NOW())`
	queryUpdate := `update "user" set coins=coins-$1 where id=$2;`
	tx, err := m.db.Begin(ctx)
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

func (m *Merch) GetByName(ctx context.Context, name string) (*entity.Merch, error) {
	query := `select id, name, cost from merch where name=$1`
	var res entity.Merch
	err := m.db.QueryRow(ctx, query, name).Scan(&res.ID, &res.Name, &res.Cost)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return nil, nil
		}
		return nil, err
	}
	return &res, nil
}

func (m *Merch) GetInventoryHistory(ctx context.Context, id uint32) ([]entity.Inventory, error) {
	query := `select m.name, count(i.merch_id) as quantity from inventory as i
				JOIN merch as m ON i.merch_id=m.id 
				WHERE i.user_id=$1
				GROUP BY m.name;`
	res := []entity.Inventory{}
	rows, err := m.db.Query(ctx, query, id)
	if err != nil {
		return res, err
	}
	for rows.Next() {
		var i entity.Inventory
		err := rows.Scan(&i.Type, &i.Quantity)
		if err != nil {
			return []entity.Inventory{}, err
		}
		res = append(res, i)
	}
	return res, nil
}
