package usecase

import (
	"avito-winter-2025/internal/entity"
	"avito-winter-2025/internal/repo/mock"
	myErrors "avito-winter-2025/internal/utils/errors"
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type args struct {
	userId    uint32
	merchName string
	merchCost uint32
	merchId   uint32
}

var ErrDB = errors.New("db error")

func TestMerchUsecase_Buy(t *testing.T) {
	tests := []struct {
		name     string
		repoMock func(context context.Context, merchRepo *mock.MockMerchInterface, coinRepo *mock.MockCoinInterface, name string, userId, merchId, cost uint32)
		args     args
		want     error
	}{
		{
			name: "Err in Get merch by name",
			repoMock: func(ctx context.Context, merchRepo *mock.MockMerchInterface, coinRepo *mock.MockCoinInterface, name string, userId, merchId, cost uint32) {
				merchRepo.EXPECT().GetByName(ctx, name).Return(nil, ErrDB)
			},
			args: args{
				userId:    1,
				merchName: "t-shirt",
			},
			want: ErrDB,
		},
		{
			name: "There are no merch with this name",
			repoMock: func(ctx context.Context, merchRepo *mock.MockMerchInterface, coinRepo *mock.MockCoinInterface, name string, userId, merchId, cost uint32) {
				merchRepo.EXPECT().GetByName(ctx, name).Return(nil, nil)
			},
			args: args{
				userId:    1,
				merchName: "T-shi",
			},
			want: myErrors.NoMerchErr,
		},
		{
			name: "Err in CheckBalance",
			repoMock: func(ctx context.Context, merchRepo *mock.MockMerchInterface, coinRepo *mock.MockCoinInterface, name string, userId, merchId, cost uint32) {
				merchRepo.EXPECT().GetByName(ctx, name).
					Return(&entity.Merch{
						ID:   merchId,
						Name: name,
						Cost: cost,
					}, nil)
				coinRepo.EXPECT().CheckBalance(ctx, userId).Return(uint32(0), ErrDB)
			},
			args: args{
				userId:    1,
				merchName: "T-shirt",
				merchId:   3,
				merchCost: 80,
			},
			want: ErrDB,
		},
		{
			name: "Err Not Enough money",
			repoMock: func(ctx context.Context, merchRepo *mock.MockMerchInterface, coinRepo *mock.MockCoinInterface, name string, userId, merchId, cost uint32) {
				merchRepo.EXPECT().GetByName(ctx, name).
					Return(&entity.Merch{
						ID:   merchId,
						Name: name,
						Cost: cost,
					}, nil)
				coinRepo.EXPECT().CheckBalance(ctx, userId).Return(uint32(60), myErrors.NotEnoughCoinErr)
			},
			args: args{
				userId:    1,
				merchName: "T-shirt",
				merchId:   3,
				merchCost: 80,
			},
			want: myErrors.NotEnoughCoinErr,
		},
		{
			name: "Err in Buying merch",
			repoMock: func(ctx context.Context, merchRepo *mock.MockMerchInterface, coinRepo *mock.MockCoinInterface, name string, userId, merchId, cost uint32) {
				merchRepo.EXPECT().GetByName(ctx, name).
					Return(&entity.Merch{
						ID:   merchId,
						Name: name,
						Cost: cost,
					}, nil)
				coinRepo.EXPECT().CheckBalance(ctx, userId).Return(uint32(100), nil)
				merchRepo.EXPECT().Buy(ctx, userId, merchId, cost).Return(ErrDB)
			},
			args: args{
				userId:    1,
				merchName: "T-shirt",
				merchId:   3,
				merchCost: 80,
			},
			want: ErrDB,
		},
		{
			name: "Success",
			repoMock: func(ctx context.Context, merchRepo *mock.MockMerchInterface, coinRepo *mock.MockCoinInterface, name string, userId, merchId, cost uint32) {
				merchRepo.EXPECT().GetByName(ctx, name).
					Return(&entity.Merch{
						ID:   merchId,
						Name: name,
						Cost: cost,
					}, nil)
				coinRepo.EXPECT().CheckBalance(ctx, userId).Return(uint32(100), nil)
				merchRepo.EXPECT().Buy(ctx, userId, merchId, cost).Return(nil)
			},
			args: args{
				userId:    1,
				merchName: "T-shirt",
				merchId:   3,
				merchCost: 80,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			coinRepo := mock.NewMockCoinInterface(ctl)
			merchRepo := mock.NewMockMerchInterface(ctl)
			usecase := NewMerch(merchRepo, coinRepo)

			tt.repoMock(context.Background(), merchRepo, coinRepo, tt.args.merchName, tt.args.userId, tt.args.merchId, tt.args.merchCost)
			got := usecase.Buy(context.Background(), tt.args.userId, tt.args.merchName)

			if !assert.Equal(t, got, tt.want) {
				t.Errorf("MerchUsecase.Buy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMerchUsecase_GetInventoryHistory(t *testing.T) {
	tests := []struct {
		name      string
		repoMock  func(ctx context.Context, merchRepo *mock.MockMerchInterface, coinRepo *mock.MockCoinInterface, id uint32)
		id        uint32
		wantError bool
		err       error
		want      []entity.Inventory
	}{
		{
			name: "Fail",
			repoMock: func(ctx context.Context, merchRepo *mock.MockMerchInterface, coinRepo *mock.MockCoinInterface, id uint32) {
				merchRepo.EXPECT().GetInventoryHistory(ctx, id).Return(nil, ErrDB)
			},
			wantError: true,
			err:       ErrDB,
			want:      []entity.Inventory{},
		},
		{
			name: "Success",
			repoMock: func(ctx context.Context, merchRepo *mock.MockMerchInterface, coinRepo *mock.MockCoinInterface, id uint32) {
				merchRepo.EXPECT().GetInventoryHistory(ctx, id).
					Return([]entity.Inventory{
						{Type: "Pen", Quantity: 2},
						{Type: "T-shirt", Quantity: 3},
					}, nil)
			},
			wantError: false,
			err:       ErrDB,
			want: []entity.Inventory{
				{Type: "Pen", Quantity: 2},
				{Type: "T-shirt", Quantity: 3},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			coinRepo := mock.NewMockCoinInterface(ctl)
			merchRepo := mock.NewMockMerchInterface(ctl)
			usecase := NewMerch(merchRepo, coinRepo)

			tt.repoMock(context.Background(), merchRepo, coinRepo, tt.id)
			got, err := usecase.GetInventoryHistory(context.Background(), tt.id)

			if (err != nil) != tt.wantError {
				t.Errorf("NoteUsecase.GetAllNotes() error = %v, wantErr %v", err, tt.wantError)
				return
			}
			if !assert.Equal(t, got, tt.want) {
				t.Errorf("MerchUsecase.GetInventory() = %v, want %v", got, tt.want)
			}
		})
	}
}
