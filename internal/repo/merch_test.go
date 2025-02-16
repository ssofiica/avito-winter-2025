package repo

import (
	"avito-winter-2025/internal/entity"
	"context"
	"testing"

	"github.com/jackc/pgx"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

func TestMerch_GetByName(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewMerch(mock)
	query := `select id, name, cost from merch where name=\$1`

	tests := []struct {
		name      string
		merchName string
		mock      func(mock pgxmock.PgxPoolIface, query string, name string)
		want      *entity.Merch
		err       error
	}{
		{
			name:      "Success",
			merchName: "pen",
			mock: func(m pgxmock.PgxPoolIface, query string, name string) {
				m.ExpectQuery(query).WithArgs(name).
					WillReturnRows(pgxmock.NewRows([]string{"id", "name", "cost"}).
						AddRow(uint32(1), "pen", uint32(100)))
			},
			want: &entity.Merch{ID: 1, Name: "pen", Cost: 100},
			err:  nil,
		},
		{
			name:      "Fail, db error",
			merchName: "pen",
			mock: func(m pgxmock.PgxPoolIface, query string, name string) {
				m.ExpectQuery(query).WithArgs(name).
					WillReturnError(ErrDB)
			},
			want: nil,
			err:  ErrDB,
		},
		{
			name:      "Fail, db error",
			merchName: "p",
			mock: func(m pgxmock.PgxPoolIface, query string, name string) {
				m.ExpectQuery(query).WithArgs(name).
					WillReturnError(pgx.ErrNoRows)
			},
			want: nil,
			err:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(mock, query, tt.merchName)
			res, err := repo.GetByName(context.Background(), tt.merchName)

			assert.Equal(t, tt.want, res)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestMerch_GetInventoryHistory(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewMerch(mock)
	query := `select m.name, count\(i.merch_id\) as quantity from inventory as i
	JOIN merch as m ON i.merch_id=m.id
	WHERE i.user_id=\$1
	GROUP BY m.name;`
	id := uint32(1)

	tests := []struct {
		name string
		mock func(mock pgxmock.PgxPoolIface, query string, id uint32)
		err  error
		want []entity.Inventory
	}{
		{
			name: "Success",
			mock: func(m pgxmock.PgxPoolIface, query string, id uint32) {
				m.ExpectQuery(query).WithArgs(id).WillReturnRows(
					pgxmock.NewRows([]string{"name", "quantity"}).
						AddRow("pen", uint32(3)).
						AddRow("t-shirt", uint32(2)))
			},
			err: nil,
			want: []entity.Inventory{
				{Type: "pen", Quantity: 3},
				{Type: "t-shirt", Quantity: 2},
			},
		},
		{
			name: "Fail",
			mock: func(m pgxmock.PgxPoolIface, query string, id uint32) {
				m.ExpectQuery(query).WithArgs(id).
					WillReturnError(ErrDB)
			},
			err:  ErrDB,
			want: []entity.Inventory{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(mock, query, id)
			res, err := repo.GetInventoryHistory(context.Background(), id)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.want, res)
		})
	}
}
