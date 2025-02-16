package repo

import (
	"avito-winter-2025/internal/entity"
	"context"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

func TestCoin_CheckBalance(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewCoin(mock)
	id := uint32(1)
	query := `select coins from "user" where id=\$1;`

	tests := []struct {
		name string
		mock func(mock pgxmock.PgxPoolIface, query string, id uint32)
		want uint32
		err  error
	}{
		{
			name: "Success",
			mock: func(m pgxmock.PgxPoolIface, query string, id uint32) {
				m.ExpectQuery(query).WithArgs(id).
					WillReturnRows(pgxmock.NewRows([]string{"coins"}).AddRow(uint32(1000)))
			},
			want: 1000,
			err:  nil,
		},
		{
			name: "Fail",
			mock: func(m pgxmock.PgxPoolIface, query string, id uint32) {
				m.ExpectQuery(query).WithArgs(id).
					WillReturnError(ErrDB)
			},
			want: 0,
			err:  ErrDB,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(mock, query, id)
			res, err := repo.CheckBalance(context.Background(), id)

			assert.Equal(t, tt.want, res)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestCoin_GetCoinHistory(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewCoin(mock)
	query := `select from_user, to_user, amount from coin_history where from_user=\$1 OR to_user=\$1;`
	id := uint32(1)

	tests := []struct {
		name string
		mock func(mock pgxmock.PgxPoolIface, query string, id uint32)
		err  error
		want []entity.Transaction
	}{
		{
			name: "Success",
			mock: func(m pgxmock.PgxPoolIface, query string, id uint32) {
				m.ExpectQuery(query).WithArgs(id).WillReturnRows(
					pgxmock.NewRows([]string{"from_user", "to_user", "amount"}).
						AddRow(uint32(1), uint32(2), uint32(100)).
						AddRow(uint32(3), uint32(1), uint32(50)))
			},
			err: nil,
			want: []entity.Transaction{
				{From: 1, To: 2, Amount: 100},
				{From: 3, To: 1, Amount: 50},
			},
		},
		{
			name: "Success, but empty",
			mock: func(m pgxmock.PgxPoolIface, query string, id uint32) {
				m.ExpectQuery(query).WithArgs(id).WillReturnRows(
					pgxmock.NewRows([]string{"from_user", "to_user", "amount"}))
			},
			err:  nil,
			want: []entity.Transaction{},
		},
		{
			name: "Fail",
			mock: func(m pgxmock.PgxPoolIface, query string, id uint32) {
				m.ExpectQuery(query).WithArgs(id).WillReturnError(ErrDB)
			},
			err:  ErrDB,
			want: []entity.Transaction{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(mock, query, id)
			res, err := repo.GetCoinHistory(context.Background(), id)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.want, res)
		})
	}
}
