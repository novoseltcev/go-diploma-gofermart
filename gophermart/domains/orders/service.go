package orders

import (
	"context"
	"errors"

	"github.com/novoseltcev/go-diploma-gofermart/gophermart/models"
	"github.com/novoseltcev/go-diploma-gofermart/gophermart/utils"
)


type OrderStorager interface {
	GetUserOrders(ctx context.Context, userID uint64) ([]models.Order, error)
	GetByNumber(ctx context.Context, number string) (*models.Order, error)
	Create(ctx context.Context, userID uint64, number string) error
}

var (
	ErrLunhNumberValidation = errors.New("invalid Lunh's number")
	ErrOrderLoaded = errors.New("order already loaded")
	ErrOrderNotMeLoaded = errors.New("order already loaded, but not owned by this user")
)


func GetUserOrders(ctx context.Context, storager OrderStorager, userID uint64) ([]models.Order, error) {
	return storager.GetUserOrders(ctx, userID)
}

func AddOrderToUser(ctx context.Context, storager OrderStorager, userID uint64, number string) error {
	if !utils.ValidateLunhNumber(number) {
		return ErrLunhNumberValidation
	}
	
	order, err := storager.GetByNumber(ctx, number)
	if err != nil {
		return err
	}

	if order != nil {
		if order.UserID == userID {
			return ErrOrderLoaded
		}
		return ErrOrderNotMeLoaded
	}

	return storager.Create(ctx, userID, number)
}
