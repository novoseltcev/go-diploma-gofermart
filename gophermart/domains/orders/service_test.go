package orders

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/novoseltcev/go-diploma-gofermart/gophermart/models"
	mock_orders "github.com/novoseltcev/go-diploma-gofermart/mocks/gophermart/domains/orders"
)

var (
	CTX = context.TODO()
	USER_ID uint64 = 1
	NUMBER = "1384858"
	ERR = errors.New("")
)

func TestSuccessGetUserOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	storager := mock_orders.NewMockOrderStorager(ctrl)

	want := []models.Order{
		{
			Number:     "",
			Status:     "",
			Accrual:    nil,
			UserId:     USER_ID,
			UploadedAt: time.Now(),
		},
	}

	storager.EXPECT().GetUserOrders(CTX, USER_ID).Return(want, nil)

	got, err := GetUserOrders(CTX, storager, USER_ID)
	require.Nil(t, err)
	assert.EqualValues(t, want, got)
}

func TestFailedGetUserOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	storager := mock_orders.NewMockOrderStorager(ctrl)

	storager.EXPECT().GetUserOrders(CTX, USER_ID).Return(nil, ERR)

	_, err := GetUserOrders(CTX, storager, USER_ID)
	assert.ErrorIs(t, err, ERR)
}

func TestSuccessAddOrderToUser(t *testing.T) {
	tests := []struct {
		name   string
		number string
	}{
		{
			name: "base",
			number: NUMBER,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			storager := mock_orders.NewMockOrderStorager(ctrl)
			
			gomock.InOrder(
				storager.EXPECT().GetByNumber(CTX, tt.number).Return(nil, nil),
				storager.EXPECT().Create(CTX, USER_ID, tt.number).Return(nil),
			)

			assert.Nil(t, AddOrderToUser(CTX, storager, USER_ID, tt.number))
		})
	}
}

func TestAddAlreadyAdderOrderToUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	storager := mock_orders.NewMockOrderStorager(ctrl)

	storager.EXPECT().GetByNumber(CTX, NUMBER).Return(&models.Order{UserId: USER_ID}, nil)
	storager.EXPECT().Create(CTX, USER_ID, NUMBER).Times(0)

	err := AddOrderToUser(CTX, storager, USER_ID, NUMBER)
	require.NotNil(t, err)
	assert.ErrorIs(t, err, OrderLoadedErr)
}

func TestAddAlreadyAdderOrderToAnotherUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	storager := mock_orders.NewMockOrderStorager(ctrl)

	storager.EXPECT().GetByNumber(CTX, NUMBER).Return(&models.Order{UserId: USER_ID + 1}, nil)
	storager.EXPECT().Create(CTX, USER_ID, NUMBER).Times(0)

	err := AddOrderToUser(CTX, storager, USER_ID, NUMBER)
	require.NotNil(t, err)
	assert.ErrorIs(t, err, OrderNotMeLoadedErr)
}

func TestFailedAddInvalidOrderToUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	storager := mock_orders.NewMockOrderStorager(ctrl)

	err := AddOrderToUser(CTX, storager, USER_ID, "21321")
	assert.ErrorIs(t, err, LunhNumberValidationErr)
}

func TestFailedAddOrderToUser(t *testing.T) {
	tests := []struct {
		name   string
		GetByNumberErr error
		CreateErr error
	}{
		{
			name: "failed GetByNumber",
			GetByNumberErr: ERR,
			CreateErr: nil,
		},
		{
			name: "failed Create",
			GetByNumberErr: nil,
			CreateErr: ERR,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			storager := mock_orders.NewMockOrderStorager(ctrl)


			storager.EXPECT().GetByNumber(CTX, NUMBER).Return(nil, tt.GetByNumberErr)
			if tt.GetByNumberErr == nil {
				require.NotNil(t, tt.CreateErr)
				storager.EXPECT().Create(CTX, USER_ID, NUMBER).Return(tt.CreateErr)
			}
			assert.ErrorIs(t, AddOrderToUser(CTX, storager, USER_ID, NUMBER), ERR)
		})
	}
}
