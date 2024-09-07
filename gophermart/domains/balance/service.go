package balance

import (
	"context"
	"errors"

	"github.com/novoseltcev/go-diploma-gofermart/gophermart/models"
	"github.com/novoseltcev/go-diploma-gofermart/gophermart/utils"
)


type BalanceStorager interface {
	GetBalance(ctx context.Context, userID uint64) (models.Money, error)
	UpdateBalance(ctx context.Context, userID uint64, value float32) error
	GetTotalWithrawn(ctx context.Context, userID uint64) (models.Money, error)
	GetUserWithdrawals(ctx context.Context, userID uint64) (result []models.Withdraw, err error)
	CreateWithdrawal(ctx context.Context, userID uint64, sum uint64, order string) error
}


var (
	ErrLunhNumberValidation = errors.New("invalid Lunh's number")
	ErrNotEnought = errors.New("not enouht balance to withdraw")
)

func GetBalance(ctx context.Context, storager BalanceStorager, userID uint64) (*models.Balance, error) {
	balance, err := storager.GetBalance(ctx, userID)
	if err != nil {
		return nil, err
	}

	withdrawn, err := storager.GetTotalWithrawn(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &models.Balance{
		Balance: balance.Value,
		Withdrawn: withdrawn.Value,
	}, nil
}

func GetUserWithdrawals(ctx context.Context, storager BalanceStorager, userID uint64) ([]models.Withdraw, error) {
	return storager.GetUserWithdrawals(ctx, userID)
}

func Withdrawn(ctx context.Context, storager BalanceStorager, userID uint64, sum uint64, order string) error {
	if !utils.ValidateLunhNumber(order) {
		return ErrLunhNumberValidation
	}

	balance, err := storager.GetBalance(ctx, userID)
	if err != nil {
		return err
	}

	newBalanceValue := balance.Value - float32(sum)
	if newBalanceValue < 0. {
		return ErrNotEnought
	}

	if err := storager.CreateWithdrawal(ctx, userID, sum, order); err != nil {
		return err
	}
	return storager.UpdateBalance(ctx, userID, newBalanceValue)
}
