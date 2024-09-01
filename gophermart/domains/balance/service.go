package balance

import (
	"context"
	"errors"

	"github.com/novoseltcev/go-diploma-gofermart/gophermart/models"
	"github.com/novoseltcev/go-diploma-gofermart/gophermart/utils"
)


type BalanceStorager interface {
	GetBalance(ctx context.Context, userId uint64) (models.Money, error)
	UpdateBalance(ctx context.Context, userId uint64, value float32) error
	GetTotalWithrawn(ctx context.Context, userId uint64) (models.Money, error)
	GetUserWithdrawals(ctx context.Context, userId uint64) (result []models.Withdraw, err error)
	CreateWithdrawal(ctx context.Context, userId uint64, sum uint64, order string) error
}


var (
	LunhNumberValidationErr = errors.New("Invalid Lunh's number")
	NotEnoughtErr = errors.New("Not enouht balance to withdraw")
)

func GetBalance(ctx context.Context, storager BalanceStorager, userId uint64) (*models.Balance, error) {
	balance, err := storager.GetBalance(ctx, userId)
	if err != nil {
		return nil, err
	}

	withdrawn, err := storager.GetTotalWithrawn(ctx, userId)
	if err != nil {
		return nil, err
	}

	return &models.Balance{
		Balance: balance.Value,
		Withdrawn: withdrawn.Value,
	}, nil
}

func GetUserWithdrawals(ctx context.Context, storager BalanceStorager, userId uint64) ([]models.Withdraw, error) {
	return storager.GetUserWithdrawals(ctx, userId)
}

func Withdrawn(ctx context.Context, storager BalanceStorager, userId uint64, sum uint64, order string) error {
	if !utils.ValidateLunhNumber(order) {
		return LunhNumberValidationErr
	}

	balance, err := storager.GetBalance(ctx, userId)
	if err != nil {
		return err
	}

	newBalanceValue := balance.Value - float32(sum)
	if newBalanceValue < 0. {
		return NotEnoughtErr
	}

	if err := storager.CreateWithdrawal(ctx, userId, sum, order); err != nil {
		return err
	}
	return storager.UpdateBalance(ctx, userId, newBalanceValue)
}
