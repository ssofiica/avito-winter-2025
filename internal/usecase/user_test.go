package usecase

import (
	"avito-winter-2025/internal/entity"
	"avito-winter-2025/internal/repo/mock"
	myErrors "avito-winter-2025/internal/utils/errors"
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUserUsecase_GetUser(t *testing.T) {
	id := uint32(1)
	name := ""
	tests := []struct {
		name      string
		repoMock  func(ctx context.Context, userRepo *mock.MockUserInterface)
		id        uint32
		userName  string
		wantError bool
		err       error
		want      entity.User
	}{
		{
			name: "Err GetUser",
			repoMock: func(ctx context.Context, userRepo *mock.MockUserInterface) {
				userRepo.EXPECT().GetUser(ctx, name, id).Return(nil, ErrDB)
			},
			wantError: true,
			err:       ErrDB,
			want:      entity.User{},
		},
		{
			name: "Success, but user is empty",
			repoMock: func(ctx context.Context, userRepo *mock.MockUserInterface) {
				userRepo.EXPECT().GetUser(ctx, name, id).Return(nil, nil)
			},
			wantError: true,
			err:       myErrors.NoUserErr,
			want:      entity.User{},
		},
		{
			name: "Success",
			repoMock: func(ctx context.Context, userRepo *mock.MockUserInterface) {
				userRepo.EXPECT().GetUser(ctx, name, id).Return(&entity.User{ID: 1, Name: "sofia", Coins: 100}, nil)
			},
			wantError: false,
			err:       nil,
			want:      entity.User{ID: 1, Name: "sofia", Coins: 100},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			userRepo := mock.NewMockUserInterface(ctl)
			usecase := NewUser(userRepo)

			tt.repoMock(context.Background(), userRepo)
			got, err := usecase.GetUser(context.Background(), name, id)

			if (err != nil) != tt.wantError {
				t.Errorf("UserUsecase.GetUser() error = %v, wantErr %v", err, tt.wantError)
				return
			}
			if !assert.Equal(t, got, tt.want) {
				t.Errorf("UserUsecase.GetUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserUsecase_Auth(t *testing.T) {
	data := entity.AuthRequest{Name: "mary", Password: "12345678M"}
	tests := []struct {
		name      string
		repoMock  func(ctx context.Context, userRepo *mock.MockUserInterface, data entity.AuthRequest)
		args      entity.AuthRequest
		wantError bool
		err       error
		want      entity.User
	}{
		{
			name: "Err GetUser",
			repoMock: func(ctx context.Context, userRepo *mock.MockUserInterface, data entity.AuthRequest) {
				userRepo.EXPECT().GetUser(ctx, data.Name, uint32(0)).Return(nil, ErrDB)
			},
			wantError: true,
			err:       ErrDB,
			want:      entity.User{},
		},
		{
			name: "Registration, err in CreateUser",
			repoMock: func(ctx context.Context, userRepo *mock.MockUserInterface, data entity.AuthRequest) {
				userRepo.EXPECT().GetUser(ctx, data.Name, uint32(0)).Return(nil, nil)
				userRepo.EXPECT().CreateUser(ctx, data.Name, gomock.Any()).Return(entity.User{}, ErrDB)
			},
			wantError: true,
			err:       myErrors.NoUserErr,
			want:      entity.User{},
		},
		{
			name: "Success Registration",
			repoMock: func(ctx context.Context, userRepo *mock.MockUserInterface, data entity.AuthRequest) {
				userRepo.EXPECT().GetUser(ctx, data.Name, uint32(0)).Return(nil, nil)
				userRepo.EXPECT().CreateUser(ctx, data.Name, gomock.Any()).Return(entity.User{ID: 1, Name: "mary", Coins: 1000}, nil)
			},
			wantError: false,
			err:       nil,
			want:      entity.User{ID: 1, Name: "mary", Coins: 1000},
		},
		{
			name: "Auth, err in Get Password",
			repoMock: func(ctx context.Context, userRepo *mock.MockUserInterface, data entity.AuthRequest) {
				userRepo.EXPECT().GetUser(ctx, data.Name, uint32(0)).Return(&entity.User{ID: 1, Name: "mary", Coins: 1000}, nil)
				userRepo.EXPECT().GetPassword(ctx, uint32(1)).Return(entity.Password(""), ErrDB)
			},
			wantError: true,
			err:       ErrDB,
			want:      entity.User{},
		},
		{
			name: "Auth, not equal passwords",
			repoMock: func(ctx context.Context, userRepo *mock.MockUserInterface, data entity.AuthRequest) {
				anotherPass := entity.Password("$2a$10$UxLlaLi4rOeWHGSDShFRD.Jtaw6wfjwYfXvlqLXDx7XihxajhdPHa")
				userRepo.EXPECT().GetUser(ctx, data.Name, uint32(0)).Return(&entity.User{ID: 1, Name: "mary", Coins: 1000}, nil)
				userRepo.EXPECT().GetPassword(ctx, uint32(1)).Return(anotherPass, nil)
			},
			wantError: true,
			err:       myErrors.WrongLoginOrPasswordErr,
			want:      entity.User{},
		},
		{
			name: "Success Auth",
			repoMock: func(ctx context.Context, userRepo *mock.MockUserInterface, data entity.AuthRequest) {
				rightPassword := entity.Password("$2a$10$lrYN1.0L/5NOcDHawDxJpOtn4jouB53uouoz8WnGFCUUDtY97Li/G")
				userRepo.EXPECT().GetUser(ctx, data.Name, uint32(0)).Return(&entity.User{ID: 1, Name: "mary", Coins: 1000}, nil)
				userRepo.EXPECT().GetPassword(ctx, uint32(1)).Return(rightPassword, nil)
			},
			wantError: false,
			err:       nil,
			want:      entity.User{ID: 1, Name: "mary", Coins: 1000},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			userRepo := mock.NewMockUserInterface(ctl)
			usecase := NewUser(userRepo)

			tt.repoMock(context.Background(), userRepo, data)
			got, err := usecase.Auth(context.Background(), data)

			if (err != nil) != tt.wantError {
				t.Errorf("UserUsecase.Auth() error = %v, wantErr %v", err, tt.wantError)
				return
			}
			if !assert.Equal(t, got, tt.want) {
				t.Errorf("UserUsecase.Auth() = %v, want %v", got, tt.want)
			}
		})
	}
}
