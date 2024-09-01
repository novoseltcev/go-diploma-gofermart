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
	CTX = context.TODO()
	USER_ID uint64 = 1
	NUMBER = "1384858"
	ERR = errors.New("")
)

func TestSuccessGetBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	storager := mock_balance.NewMockBalanceStorager(ctrl)

	want := models.Balance{
		Balance: 21.23232,
		Withdrawn: 12.,
	}
	
	gomock.InOrder(
		storager.EXPECT().GetBalance(CTX, USER_ID).Return(models.Money{Value: want.Balance}, nil),
		storager.EXPECT().GetTotalWithrawn(CTX, USER_ID).Return(models.Money{Value: want.Withdrawn}, nil),
	)

	got, err := GetBalance(CTX, storager, USER_ID)
	require.Nil(t, err)
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
			Sum: models.Money{Value: 1},
		},
		{
			Order: "4435464",
			Sum: models.Money{Value: 2565.},
		},
	}

	storager.EXPECT().GetUserWithdrawals(CTX, USER_ID).Return(want, nil)

	got, err := GetUserWithdrawals(CTX, storager, USER_ID)
	require.Nil(t, err)
	assert.EqualValues(t, want, got)
}

func TestFailedGetUserWithdrawals(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	storager := mock_balance.NewMockBalanceStorager(ctrl)

	storager.EXPECT().GetUserWithdrawals(CTX, USER_ID).Return(nil, ERR)

	_, err := GetUserWithdrawals(CTX, storager, USER_ID)
	assert.ErrorIs(t, err, ERR)
}

func TestSuccessWithdrawn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	storager := mock_balance.NewMockBalanceStorager(ctrl)
	var sum uint64 = 123

	gomock.InOrder(
		storager.EXPECT().GetBalance(CTX, USER_ID).Return(models.Money{Value: float32(sum)}, nil),
		storager.EXPECT().CreateWithdrawal(CTX, USER_ID, sum, NUMBER).Return(nil),
		storager.EXPECT().UpdateBalance(CTX, USER_ID, float32(0)).Return(nil),
	)

	assert.Nil(t, Withdrawn(CTX, storager, USER_ID, sum, NUMBER))
}

func TestFailedInvalidNumberWithdrawn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	storager := mock_balance.NewMockBalanceStorager(ctrl)
	var sum uint64 = 123

	err := Withdrawn(CTX, storager, USER_ID, sum, "21321")
	assert.ErrorIs(t, err, LunhNumberValidationErr)
}

func TestFailedNotEnoughtWithdrawn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	storager := mock_balance.NewMockBalanceStorager(ctrl)
	var sum uint64 = 123

	storager.EXPECT().GetBalance(CTX, USER_ID).Return(models.Money{Value: float32(122.99)}, nil)
	storager.EXPECT().CreateWithdrawal(CTX, USER_ID, sum, NUMBER).Times(0)
	storager.EXPECT().UpdateBalance(CTX, USER_ID, gomock.Any()).Times(0)

	err := Withdrawn(CTX, storager, USER_ID, sum, NUMBER)
	assert.ErrorIs(t, err, NotEnoughtErr)
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
			GetBalanceErr: ERR,
			CreateWithdrawalErr: nil,
			UpdateBalanceErr: nil,
		},
		{
			name: "failed CreateWithdrawal",
			GetBalanceErr: nil,
			CreateWithdrawalErr: ERR,
			UpdateBalanceErr: nil,
		},
		{
			name: "failed UpdateBalance",
			GetBalanceErr: nil,
			CreateWithdrawalErr: nil,
			UpdateBalanceErr: ERR,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			storager := mock_balance.NewMockBalanceStorager(ctrl)
			var sum uint64 = 1

			storager.EXPECT().GetBalance(CTX, USER_ID).Return(models.Money{Value: float32(sum)}, tt.GetBalanceErr)
			if tt.GetBalanceErr == nil {
				storager.EXPECT().CreateWithdrawal(CTX, USER_ID, sum, NUMBER).Return(tt.CreateWithdrawalErr)
				if tt.CreateWithdrawalErr == nil {
					storager.EXPECT().UpdateBalance(CTX, USER_ID, float32(0)).Return(tt.UpdateBalanceErr)
				}
			}
			assert.ErrorIs(t, Withdrawn(CTX, storager, USER_ID, sum, NUMBER), ERR)
		})
	}
}
