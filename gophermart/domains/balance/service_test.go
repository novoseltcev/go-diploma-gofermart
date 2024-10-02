package balance

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/novoseltcev/go-diploma-gofermart/gophermart/models"
	mock_balance "github.com/novoseltcev/go-diploma-gofermart/mocks/gophermart/domains/balance"
)

var (
	someCtx = context.TODO()
	someUserID uint64 = 1
	someNumber = "1384858"
	errSome = errors.New("")
)

func TestSuccessGetBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	storager := mock_balance.NewMockBalanceStorager(ctrl)

	want := models.Balance{
		Balance: 123.45,
		Withdrawn: 12.,
	}
	
	gomock.InOrder(
		storager.EXPECT().GetBalance(someCtx, someUserID).Return(want.Balance, nil),
		storager.EXPECT().GetTotalWithrawn(someCtx, someUserID).Return(want.Withdrawn, nil),
	)

	got, someErr := GetBalance(someCtx, storager, someUserID)
	require.Nil(t, someErr)
	require.NotNil(t, got)
	assert.EqualValues(t, want, *got)
}


func TestSuccessGetUserWithdrawals(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	storager := mock_balance.NewMockBalanceStorager(ctrl)

	want := []models.Withdraw{
		{
			Order: "123213",
			Sum: 1,
		},
		{
			Order: "4435464",
			Sum: 2565.,
		},
	}

	storager.EXPECT().GetUserWithdrawals(someCtx, someUserID).Return(want, nil)

	got, someErr := GetUserWithdrawals(someCtx, storager, someUserID)
	require.Nil(t, someErr)
	assert.EqualValues(t, want, got)
}

func TestFailedGetUserWithdrawals(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	storager := mock_balance.NewMockBalanceStorager(ctrl)

	storager.EXPECT().GetUserWithdrawals(someCtx, someUserID).Return(nil, errSome)

	_, someErr := GetUserWithdrawals(someCtx, storager, someUserID)
	assert.ErrorIs(t, someErr, someErr)
}

func TestSuccessWithdrawn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	storager := mock_balance.NewMockBalanceStorager(ctrl)
	var sum models.Money = 123.12

	gomock.InOrder(
		storager.EXPECT().GetBalance(someCtx, someUserID).Return(sum, nil),
		storager.EXPECT().CreateWithdrawal(someCtx, someUserID, sum, someNumber).Return(nil),
		storager.EXPECT().UpdateBalance(someCtx, someUserID, float32(0)).Return(nil),
	)

	assert.Nil(t, Withdrawn(someCtx, storager, someUserID, sum, someNumber))
}

func TestFailedInvalidNumberWithdrawn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	storager := mock_balance.NewMockBalanceStorager(ctrl)
	var sum models.Money = 123.12

	someErr := Withdrawn(someCtx, storager, someUserID, sum, "21321")
	assert.ErrorIs(t, someErr, ErrLunhNumberValidation)
}

func TestFailedNotEnoughtWithdrawn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	storager := mock_balance.NewMockBalanceStorager(ctrl)
	var sum models.Money = 123.12

	storager.EXPECT().GetBalance(someCtx, someUserID).Return(122.99, nil)
	storager.EXPECT().CreateWithdrawal(someCtx, someUserID, sum, someNumber).Times(0)
	storager.EXPECT().UpdateBalance(someCtx, someUserID, gomock.Any()).Times(0)

	someErr := Withdrawn(someCtx, storager, someUserID, sum, someNumber)
	assert.ErrorIs(t, someErr, ErrNotEnought)
}

func TestFailedWithdrawn(t *testing.T) {
	tests := []struct {
		name   string
		GetBalanceErr error
		CreateWithdrawalErr error
		UpdateBalanceErr error
	}{
		{
			name: "failed GetBalance",
			GetBalanceErr: errSome,
			CreateWithdrawalErr: nil,
			UpdateBalanceErr: nil,
		},
		{
			name: "failed CreateWithdrawal",
			GetBalanceErr: nil,
			CreateWithdrawalErr: errSome,
			UpdateBalanceErr: nil,
		},
		{
			name: "failed UpdateBalance",
			GetBalanceErr: nil,
			CreateWithdrawalErr: nil,
			UpdateBalanceErr: errSome,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			storager := mock_balance.NewMockBalanceStorager(ctrl)
			var sum models.Money = 1

			storager.EXPECT().GetBalance(someCtx, someUserID).Return(sum, tt.GetBalanceErr)
			if tt.GetBalanceErr == nil {
				storager.EXPECT().CreateWithdrawal(someCtx, someUserID, sum, someNumber).Return(tt.CreateWithdrawalErr)
				if tt.CreateWithdrawalErr == nil {
					storager.EXPECT().UpdateBalance(someCtx, someUserID, float32(0)).Return(tt.UpdateBalanceErr)
				}
			}
			assert.ErrorIs(t, Withdrawn(someCtx, storager, someUserID, sum, someNumber), errSome)
		})
	}
}
