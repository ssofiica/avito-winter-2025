package repo

import (
	"avito-winter-2025/internal/entity"
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

func TestUser_GetUser(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewUser(mock)
	queryName := `select id, name, coins from "user" where name=\$1`
	queryId := `select id, name, coins from "user" where id=\$1`

	tests := []struct {
		name     string
		userName string
		userId   uint32
		query    string
		mock     func(m pgxmock.PgxPoolIface, query string, name string, id uint32)
		want     *entity.User
		err      error
	}{
		{
			name:     "Success, using name",
			userName: "sofia",
			userId:   0,
			query:    queryName,
			mock: func(m pgxmock.PgxPoolIface, query string, name string, id uint32) {
				m.ExpectQuery(query).WithArgs(name).
					WillReturnRows(pgxmock.NewRows([]string{"id", "name", "coins"}).
						AddRow(uint32(1), "sofia", uint32(1000)))
			},
			want: &entity.User{ID: 1, Name: "sofia", Coins: 1000},
			err:  nil,
		},
		{
			name:     "Success, using id",
			userName: "",
			userId:   1,
			query:    queryId,
			mock: func(m pgxmock.PgxPoolIface, query string, name string, id uint32) {
				m.ExpectQuery(query).WithArgs(id).
					WillReturnRows(pgxmock.NewRows([]string{"id", "name", "coins"}).
						AddRow(uint32(1), "sofia", uint32(1000)))
			},
			want: &entity.User{ID: 1, Name: "sofia", Coins: 1000},
			err:  nil,
		},
		{
			name:     "Fail, name and id is empty",
			userName: "",
			userId:   0,
			query:    queryId,
			mock: func(m pgxmock.PgxPoolIface, query string, name string, id uint32) {
			},
			want: nil,
			err:  errors.New("Пользователь не найден"),
		},
		{
			name:     "Fail, Err No Rows",
			userName: "",
			userId:   1,
			query:    queryId,
			mock: func(m pgxmock.PgxPoolIface, query string, name string, id uint32) {
				m.ExpectQuery(query).WithArgs(id).
					WillReturnError(pgx.ErrNoRows)
			},
			want: nil,
			err:  nil,
		},
		{
			name:     "Fail, db error",
			userName: "",
			userId:   1,
			query:    queryId,
			mock: func(m pgxmock.PgxPoolIface, query string, name string, id uint32) {
				m.ExpectQuery(query).WithArgs(id).
					WillReturnError(ErrDB)
			},
			want: nil,
			err:  ErrDB,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(mock, tt.query, tt.userName, tt.userId)
			res, err := repo.GetUser(context.Background(), tt.userName, tt.userId)

			assert.Equal(t, tt.want, res)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestUser_GetPassword(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewUser(mock)
	query := `select password from "user" where id=\$1;`
	id := uint32(1)

	tests := []struct {
		name string
		mock func(mock pgxmock.PgxPoolIface, query string, id uint32)
		err  error
		want entity.Password
	}{
		{
			name: "Success",
			mock: func(m pgxmock.PgxPoolIface, query string, id uint32) {
				m.ExpectQuery(query).WithArgs(id).WillReturnRows(
					pgxmock.NewRows([]string{"password"}).
						AddRow("1234567"))
			},
			err:  nil,
			want: entity.Password("1234567"),
		},
		{
			name: "Fail",
			mock: func(m pgxmock.PgxPoolIface, query string, id uint32) {
				m.ExpectQuery(query).WithArgs(id).
					WillReturnError(ErrDB)
			},
			err:  ErrDB,
			want: entity.Password(""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(mock, query, id)
			res, err := repo.GetPassword(context.Background(), id)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.want, res)
		})
	}
}

func TestUser_CreateUser(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewUser(mock)
	queryName := `insert into "user"\(name, password, coins\) values \(\$1, \$2, \$3\) returning id, name, coins;`
	name := "sofia"
	password := "12345"
	coins := 1000

	tests := []struct {
		name  string
		query string
		mock  func(m pgxmock.PgxPoolIface, query string)
		want  entity.User
		err   error
	}{
		{
			name:  "Success",
			query: queryName,
			mock: func(m pgxmock.PgxPoolIface, query string) {
				m.ExpectQuery(query).WithArgs(name, password, coins).
					WillReturnRows(pgxmock.NewRows([]string{"id", "name", "coins"}).
						AddRow(uint32(1), "sofia", uint32(1000)))
			},
			want: entity.User{ID: 1, Name: "sofia", Coins: 1000},
			err:  nil,
		},
		{
			name:  "Fail, this user already exists",
			query: queryName,
			mock: func(m pgxmock.PgxPoolIface, query string) {
				m.ExpectQuery(query).WithArgs(name, password, coins).
					WillReturnError(pgx.PgError{Code: "23505"})
			},
			want: entity.User{},
			err:  errors.New("Запись с указанными данными уже существует"),
		},
		{
			name:  "Fail",
			query: queryName,
			mock: func(m pgxmock.PgxPoolIface, query string) {
				m.ExpectQuery(query).WithArgs(name, password, coins).
					WillReturnError(ErrDB)
			},
			want: entity.User{},
			err:  ErrDB,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(mock, tt.query)
			res, err := repo.CreateUser(context.Background(), name, password)

			assert.Equal(t, tt.want, res)
			assert.Equal(t, tt.err, err)
		})
	}
}
