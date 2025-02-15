package repo

import (
	"avito-winter-2025/internal/entity"
	myErrors "avito-winter-2025/internal/utils/errors"
	"context"

	"github.com/jackc/pgx"
	pgx5 "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var COINS = 1000

//go:generate mockgen -source=user.go -destination=mock/user_mock.go -package=mock
type UserInterface interface {
	GetUser(ctx context.Context, name string, id uint32) (*entity.User, error)
	CreateUser(ctx context.Context, name string, password string) (entity.User, error)
	GetPassword(ctx context.Context, id uint32) (entity.Password, error)
}

type User struct {
	db *pgxpool.Pool
}

func NewUser(db *pgxpool.Pool) UserInterface {
	return &User{db: db}
}

func (u *User) GetUser(ctx context.Context, name string, id uint32) (*entity.User, error) {
	if name == "" && id == 0 {
		return nil, myErrors.NoUserErr
	}
	query1 := `select id, name, coins from "user" where `
	query2 := `=$1`
	var query string
	var res entity.User
	var row pgx5.Row
	if name != "" {
		query = query1 + `name` + query2
		row = u.db.QueryRow(ctx, query, name)
	} else if id > 0 {
		query = query1 + `id` + query2
		row = u.db.QueryRow(ctx, query, id)
	}
	err := row.Scan(&res.ID, &res.Name, &res.Coins)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return nil, nil
		}
		return nil, err
	}
	return &res, nil
}

func (u *User) CreateUser(ctx context.Context, name string, password string) (entity.User, error) {
	query := `insert into "user"(name, password, coins) values ($1, $2, $3) returning id, name, coins;`
	var res entity.User
	err := u.db.QueryRow(ctx, query, name, password, COINS).Scan(&res.ID, &res.Name, &res.Coins)
	if err != nil {
		if pgErr, ok := err.(pgx.PgError); ok {
			if pgErr.Code == "23505" {
				return entity.User{}, myErrors.NotUnique
			}
		}
		return entity.User{}, err
	}
	return res, nil
}

func (u *User) GetPassword(ctx context.Context, id uint32) (entity.Password, error) {
	query := `select password from "user" where id=$1;`
	var res string
	err := u.db.QueryRow(ctx, query, id).Scan(&res)
	if err != nil {
		return "", err
	}
	return entity.Password(res), nil
}
