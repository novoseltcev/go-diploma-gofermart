package orders

import (
	"context"
	"errors"

	"github.com/novoseltcev/go-diploma-gofermart/gophermart/models"
	"github.com/novoseltcev/go-diploma-gofermart/gophermart/utils"
)


type OrderStorager interface {
	GetUserOrders(ctx context.Context, userId uint64) ([]models.Order, error)
	GetByNumber(ctx context.Context, number string) (*models.Order, error)
	Create(ctx context.Context, userId uint64, number string) error
}

var (
	LunhNumberValidationErr = errors.New("Invalid Lunh's number")
	OrderLoadedErr = errors.New("Order already loaded")
	OrderNotMeLoadedErr = errors.New("Order already loaded, but not owned by this user")
)


func GetUserOrders(ctx context.Context, storager OrderStorager, userId uint64) ([]models.Order, error) {
	return storager.GetUserOrders(ctx, userId)
}

func AddOrderToUser(ctx context.Context, storager OrderStorager, userId uint64, number string) error {
	if !utils.ValidateLunhNumber(number) {
		return LunhNumberValidationErr
	}
	
	order, err := storager.GetByNumber(ctx, number)
	if err != nil {
		return err
	}

	if order != nil {
		if order.UserId == userId {
			return OrderLoadedErr
		}
		return OrderNotMeLoadedErr
	}

	return storager.Create(ctx, userId, number)
}
