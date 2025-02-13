package usecase

import (
	"avito-winter-2025/internal/entity"
	"avito-winter-2025/internal/repo"
	myErrors "avito-winter-2025/internal/utils/errors"
	"context"
)

type UserInterface interface {
	Auth(ctx context.Context, data entity.AuthRequest) (entity.User, error)
	GetUser(ctx context.Context, name string, id uint32) (entity.User, error)
}

type User struct {
	repo repo.UserInterface
}

func NewUser(r repo.UserInterface) UserInterface {
	return &User{repo: r}
}

func (u *User) Auth(ctx context.Context, data entity.AuthRequest) (entity.User, error) {
	user, err := u.repo.GetUser(ctx, data.Name, 0)
	if err != nil {
		return entity.User{}, err
	}
	if user == nil {
		var pass entity.Password
		err := pass.Hash(data.Password)
		if err != nil {
			return entity.User{}, err
		}
		res, err := u.repo.CreateUser(ctx, data.Name, string(pass))
		if err != nil {
			return entity.User{}, err
		}
		return res, nil
	}
	password, err := u.repo.GetPassword(ctx, user.ID)
	if err != nil {
		return entity.User{}, err
	}
	if !password.IsEqual(data.Password) {
		return entity.User{}, myErrors.WrongLoginOrPasswordErr
	}
	return *user, nil
}

func (u *User) GetUser(ctx context.Context, name string, id uint32) (entity.User, error) {
	user, err := u.repo.GetUser(ctx, name, id)
	if err != nil {
		return entity.User{}, err
	}
	if user == nil {
		return entity.User{}, myErrors.NoUserErr
	}
	return *user, nil
}
