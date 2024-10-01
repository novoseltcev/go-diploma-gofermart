package storager

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/novoseltcev/go-diploma-gofermart/gophermart/domains/orders"
	"github.com/novoseltcev/go-diploma-gofermart/gophermart/models"
	"github.com/novoseltcev/go-diploma-gofermart/shared"
)

type repository struct {
	tx *sqlx.Tx
} 

func New(uow *shared.UnitOfWork) orders.OrderStorager {
	return &repository{tx: uow.Tx}
}

func (r *repository) GetUserOrders(ctx context.Context, userID uint64) (result []models.Order, err error) {
	err = r.tx.SelectContext(ctx, &result, "SELECT number, status, accrual, user_id, uploaded_at FROM orders WHERE user_id = $1", userID)
	return result, err
}

func (r *repository) GetByNumber(ctx context.Context, number string) (*models.Order, error) {
	var order models.Order
	err := r.tx.GetContext(ctx, &order, "SELECT number, status, accrual, user_id, uploaded_at FROM orders WHERE number = $1", number)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &order, err
}

func (r *repository) Create(ctx context.Context, userID uint64, number string) error {
	_, err := r.tx.ExecContext(ctx, "INSERT INTO orders (number, user_id) VALUES ($1, $2)", number, userID)
	return err
}
