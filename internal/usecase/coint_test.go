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

type CoinArgs struct {
	From   uint32
	To     uint32
	Amount uint32
}

func TestCoinUsecase_SendCoin(t *testing.T) {
	tests := []struct {
		name     string
		repoMock func(ctx context.Context, userRepo *mock.MockUserInterface, coinRepo *mock.MockCoinInterface, from, to, amount uint32)
		args     CoinArgs
		want     error
	}{
		{
			name: "Err in CheckBalance",
			repoMock: func(ctx context.Context, userRepo *mock.MockUserInterface, coinRepo *mock.MockCoinInterface, from, to, amount uint32) {
				coinRepo.EXPECT().CheckBalance(ctx, from).Return(uint32(0), ErrDB)
			},
			args: CoinArgs{
				From:   1,
				To:     2,
				Amount: 20,
			},
			want: ErrDB,
		},
		{
			name: "Err Not enough money, balance less than amount",
			repoMock: func(ctx context.Context, userRepo *mock.MockUserInterface, coinRepo *mock.MockCoinInterface, from, to, amount uint32) {
				coinRepo.EXPECT().CheckBalance(ctx, from).Return(amount-1, nil)
			},
			args: CoinArgs{
				From:   1,
				To:     2,
				Amount: 20,
			},
			want: myErrors.NotEnoughCoinErr,
		},
		{
			name: "Err in SendCoin",
			repoMock: func(ctx context.Context, userRepo *mock.MockUserInterface, coinRepo *mock.MockCoinInterface, from, to, amount uint32) {
				coinRepo.EXPECT().CheckBalance(ctx, from).Return(amount+1, nil)
				coinRepo.EXPECT().SendCoin(ctx, entity.Transaction{
					From:   from,
					To:     to,
					Amount: amount,
				}).Return(ErrDB)
			},
			args: CoinArgs{
				From:   1,
				To:     2,
				Amount: 20,
			},
			want: ErrDB,
		},
		{
			name: "Success",
			repoMock: func(ctx context.Context, userRepo *mock.MockUserInterface, coinRepo *mock.MockCoinInterface, from, to, amount uint32) {
				coinRepo.EXPECT().CheckBalance(ctx, from).Return(amount+1, nil)
				coinRepo.EXPECT().SendCoin(ctx, entity.Transaction{
					From:   from,
					To:     to,
					Amount: amount,
				}).Return(nil)
			},
			args: CoinArgs{
				From:   1,
				To:     2,
				Amount: 20,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			coinRepo := mock.NewMockCoinInterface(ctl)
			userRepo := mock.NewMockUserInterface(ctl)
			usecase := NewCoin(coinRepo, userRepo)

			tt.repoMock(context.Background(), userRepo, coinRepo, tt.args.From, tt.args.To, tt.args.Amount)
			got := usecase.SendCoin(context.Background(), tt.args.From, tt.args.To, tt.args.Amount)

			if !assert.Equal(t, got, tt.want) {
				t.Errorf("CoinUsecase.SendCoin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCoinUsecase_GetCoinHistory(t *testing.T) {
	tests := []struct {
		name      string
		repoMock  func(ctx context.Context, userRepo *mock.MockUserInterface, coinRepo *mock.MockCoinInterface, id uint32)
		id        uint32
		wantError bool
		err       error
		want      entity.CoinHistory
	}{
		{
			name: "Err GetCoinHistory",
			repoMock: func(ctx context.Context, userRepo *mock.MockUserInterface, coinRepo *mock.MockCoinInterface, id uint32) {
				coinRepo.EXPECT().GetCoinHistory(ctx, id).Return([]entity.Transaction{}, ErrDB)
			},
			id:        1,
			wantError: true,
			err:       ErrDB,
			want:      entity.CoinHistory{Received: []entity.Received{}, Sent: []entity.Sent{}},
		},
		{
			name: "Err in GetUser",
			repoMock: func(ctx context.Context, userRepo *mock.MockUserInterface, coinRepo *mock.MockCoinInterface, id uint32) {
				coinRepo.EXPECT().GetCoinHistory(ctx, id).
					Return([]entity.Transaction{
						{From: 1, To: 2, Amount: 50},
						{From: 2, To: 1, Amount: 14},
						{From: 1, To: 3, Amount: 100},
					}, nil)
				userRepo.EXPECT().GetUser(ctx, "", uint32(2)).Return(nil, ErrDB)
			},
			id:        1,
			wantError: true,
			err:       ErrDB,
			want:      entity.CoinHistory{Received: []entity.Received{}, Sent: []entity.Sent{}},
		},
		{
			name: "Success, but CoinHistory is empty",
			repoMock: func(ctx context.Context, userRepo *mock.MockUserInterface, coinRepo *mock.MockCoinInterface, id uint32) {
				coinRepo.EXPECT().GetCoinHistory(ctx, id).
					Return([]entity.Transaction{}, nil)
			},
			id:        1,
			wantError: false,
			err:       nil,
			want:      entity.CoinHistory{Received: []entity.Received{}, Sent: []entity.Sent{}},
		},
		{
			name: "Success",
			repoMock: func(ctx context.Context, userRepo *mock.MockUserInterface, coinRepo *mock.MockCoinInterface, id uint32) {
				coinRepo.EXPECT().GetCoinHistory(ctx, id).
					Return([]entity.Transaction{
						{From: 1, To: 2, Amount: 50},
						{From: 2, To: 1, Amount: 14},
						{From: 1, To: 3, Amount: 20},
					}, nil)
				userRepo.EXPECT().GetUser(ctx, "", uint32(2)).
					Return(&entity.User{ID: 2, Name: "mary", Coins: 200}, nil)
				userRepo.EXPECT().GetUser(ctx, "", uint32(2)).
					Return(&entity.User{ID: 2, Name: "mary", Coins: 200}, nil)
				userRepo.EXPECT().GetUser(ctx, "", uint32(3)).
					Return(&entity.User{ID: 3, Name: "sofia", Coins: 100}, nil)
			},
			id:        1,
			wantError: false,
			err:       nil,
			want: entity.CoinHistory{
				Received: []entity.Received{
					{FromUser: "mary", Amount: 14},
				},
				Sent: []entity.Sent{
					{ToUser: "mary", Amount: 50},
					{ToUser: "sofia", Amount: 20},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			coinRepo := mock.NewMockCoinInterface(ctl)
			userRepo := mock.NewMockUserInterface(ctl)
			usecase := NewCoin(coinRepo, userRepo)

			tt.repoMock(context.Background(), userRepo, coinRepo, tt.id)
			got, err := usecase.GetCoinHistory(context.Background(), tt.id)

			if (err != nil) != tt.wantError {
				t.Errorf("CoinUsecase.GetCoinHistory() error = %v, wantErr %v", err, tt.wantError)
				return
			}
			if !assert.Equal(t, got, tt.want) {
				t.Errorf("CoinUsecase.GetCoinHistory() = %v, want %v", got, tt.want)
			}
		})
	}
}
