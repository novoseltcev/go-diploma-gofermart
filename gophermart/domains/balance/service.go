package balance

import (
	"context"
	"errors"

	"github.com/novoseltcev/go-diploma-gofermart/gophermart/models"
	"github.com/novoseltcev/go-diploma-gofermart/gophermart/utils"
)


type BalanceStorager interface {
	GetCurrent(ctx context.Context, userID models.UserID) (models.Money, error)
	UpdateBalance(ctx context.Context, userID models.UserID, value models.Money) error
	GetTotalWithrawn(ctx context.Context, userID models.UserID) (models.Money, error)
	GetUserWithdrawals(ctx context.Context, userID models.UserID) (result []models.Withdraw, err error)
	CreateWithdrawal(ctx context.Context, userID models.UserID, sum models.Money, order string) error
}


var (
	ErrLunhNumberValidation = errors.New("invalid Lunh's number")
	ErrNotEnought = errors.New("not enouht balance to withdraw")
)

func GetBalance(ctx context.Context, storager BalanceStorager, userID models.UserID) (*models.Balance, error) {
	current, err := storager.GetCurrent(ctx, userID)
	if err != nil {
		return nil, err
	}

	withdrawn, err := storager.GetTotalWithrawn(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &models.Balance{
		Current: current,
		Withdrawn: withdrawn,
	}, nil
}

func GetUserWithdrawals(ctx context.Context, storager BalanceStorager, userID models.UserID) ([]models.Withdraw, error) {
	return storager.GetUserWithdrawals(ctx, userID)
}

func Withdrawn(ctx context.Context, storager BalanceStorager, userID models.UserID, sum models.Money, order string) error {
	if !utils.ValidateLunhNumber(order) {
		return ErrLunhNumberValidation
	}

	current, err := storager.GetCurrent(ctx, userID)
	if err != nil {
		return err
	}

	newBalance := current - sum
	if newBalance < 0. {
		return ErrNotEnought
	}

	if err := storager.CreateWithdrawal(ctx, userID, sum, order); err != nil {
		return err
	}
	return storager.UpdateBalance(ctx, userID, newBalance)
}
