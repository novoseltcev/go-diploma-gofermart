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
	someCtx = context.TODO()
	someUserID uint64 = 1
	someNumber = "1384858"
	errSome = errors.New("")
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
			UserID:     someUserID,
			UploadedAt: time.Now(),
		},
	}

	storager.EXPECT().GetUserOrders(someCtx, someUserID).Return(want, nil)

	got, err := GetUserOrders(someCtx, storager, someUserID)
	require.Nil(t, err)
	assert.EqualValues(t, want, got)
}

func TestFailedGetUserOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	storager := mock_orders.NewMockOrderStorager(ctrl)

	storager.EXPECT().GetUserOrders(someCtx, someUserID).Return(nil, errSome)

	_, err := GetUserOrders(someCtx, storager, someUserID)
	assert.ErrorIs(t, err, errSome)
}

func TestSuccessAddOrderToUser(t *testing.T) {
	tests := []struct {
		name   string
		number string
	}{
		{
			name: "base",
			number: someNumber,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			storager := mock_orders.NewMockOrderStorager(ctrl)
			
			gomock.InOrder(
				storager.EXPECT().GetByNumber(someCtx, tt.number).Return(nil, nil),
				storager.EXPECT().Create(someCtx, someUserID, tt.number).Return(nil),
			)

			assert.Nil(t, AddOrderToUser(someCtx, storager, someUserID, tt.number))
		})
	}
}

func TestAddAlreadyAdderOrderToUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	storager := mock_orders.NewMockOrderStorager(ctrl)

	storager.EXPECT().GetByNumber(someCtx, someNumber).Return(&models.Order{UserID: someUserID}, nil)
	storager.EXPECT().Create(someCtx, someUserID, someNumber).Times(0)

	err := AddOrderToUser(someCtx, storager, someUserID, someNumber)
	require.NotNil(t, err)
	assert.ErrorIs(t, err, ErrOrderLoaded)
}

func TestAddAlreadyAdderOrderToAnotherUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	storager := mock_orders.NewMockOrderStorager(ctrl)

	storager.EXPECT().GetByNumber(someCtx, someNumber).Return(&models.Order{UserID: someUserID + 1}, nil)
	storager.EXPECT().Create(someCtx, someUserID, someNumber).Times(0)

	err := AddOrderToUser(someCtx, storager, someUserID, someNumber)
	require.NotNil(t, err)
	assert.ErrorIs(t, err, ErrOrderNotMeLoaded)
}

func TestFailedAddInvalidOrderToUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	storager := mock_orders.NewMockOrderStorager(ctrl)

	err := AddOrderToUser(someCtx, storager, someUserID, "21321")
	assert.ErrorIs(t, err, ErrLunhNumberValidation)
}

func TestFailedAddOrderToUser(t *testing.T) {
	tests := []struct {
		name   string
		GetByNumberErr error
		CreateErr error
	}{
		{
			name: "failed GetByNumber",
			GetByNumberErr: errSome,
			CreateErr: nil,
		},
		{
			name: "failed Create",
			GetByNumberErr: nil,
			CreateErr: errSome,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			storager := mock_orders.NewMockOrderStorager(ctrl)


			storager.EXPECT().GetByNumber(someCtx, someNumber).Return(nil, tt.GetByNumberErr)
			if tt.GetByNumberErr == nil {
				require.NotNil(t, tt.CreateErr)
				storager.EXPECT().Create(someCtx, someUserID, someNumber).Return(tt.CreateErr)
			}
			assert.ErrorIs(t, AddOrderToUser(someCtx, storager, someUserID, someNumber), errSome)
		})
	}
}
