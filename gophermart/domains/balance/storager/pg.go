package storager

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/novoseltcev/go-diploma-gofermart/gophermart/domains/balance"
	. "github.com/novoseltcev/go-diploma-gofermart/gophermart/models"
	"github.com/novoseltcev/go-diploma-gofermart/shared"
)

type repository struct {
	tx *sqlx.Tx
} 

func New(uow *shared.UnitOfWork) balance.BalanceStorager {
	return &repository{tx: uow.Tx}
}

func (r *repository) GetBalance(ctx context.Context, userID UserID) (result Money, err error) {
	return result, r.tx.GetContext(ctx, &result, "SELECT balance::numeric FROM users WHERE id = $1", userID)
}

func (r *repository) UpdateBalance(ctx context.Context, userID UserID, value Money) error {
	_, err := r.tx.ExecContext(ctx, "UPDATE users SET balance = $1::MONEY WHERE id = $2", value, userID)
	return err
}


func (r *repository) GetUserWithdrawals(ctx context.Context, userID UserID) (result []Withdraw, err error) {
	err = r.tx.SelectContext(ctx, &result, `SELECT "order", sum::numeric as sum FROM withdrawals WHERE user_id = $1 ORDER BY processed_at DESC`, userID)
	return result, err
}

func (r *repository) GetTotalWithrawn(ctx context.Context, userID UserID) (result Money, err error) {
	return result, r.tx.GetContext(ctx, &result, `SELECT COALESCE(SUM("sum"), 0::MONEY)::numeric as withdrawn FROM withdrawals WHERE user_id = $1`, userID)
}


func (r *repository) CreateWithdrawal(ctx context.Context, userID UserID, sum Money, order string) error {
	_, err := r.tx.ExecContext(ctx, `INSERT INTO withdrawals (user_id, sum, "order") VALUES ($1, $2::MONEY, $3)`, userID, sum, order)
	return err
}
