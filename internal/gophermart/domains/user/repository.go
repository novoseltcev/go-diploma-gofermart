package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/novoseltcev/go-diploma-gofermart/internal/gophermart/models"
	"github.com/novoseltcev/go-diploma-gofermart/internal/shared"
)

type UserRepository struct {
	tx *sqlx.Tx
}

func NewStorager(uow *shared.UnitOfWork) UserStorager {
	return &UserRepository{tx: uow.Tx}
}

func (r *UserRepository) IsLoginExists(ctx context.Context, login string) (result bool, err error) {
	err = r.tx.GetContext(ctx, &result, "SELECT COUNT(login) > 0 FROM users WHERE login = $1", login)
	return result, err
}

func (r *UserRepository) GetById(ctx context.Context, userId uint64) (*models.User, error) {
	var user models.User
	err := r.tx.GetContext(ctx, &user, "SELECT id, login, balance, hashed_password FROM users WHERE id = $1", userId)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &user, err
}

func (r *UserRepository) GetByLogin(ctx context.Context, login string) (*models.User, error) {
	var user models.User
	err := r.tx.GetContext(ctx, &user, "SELECT id, login, balance, hashed_password FROM users WHERE login = $1", login)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &user, err
}

func (r *UserRepository) Create(ctx context.Context, login, password string) (uint64, error) {
	var res uint64
	err := r.tx.GetContext(ctx, &res, "INSERT INTO users (login, hashed_password) VALUES ($1, $2) RETURNING id", login, password)
	if err != nil {
		return 0, err
	}
	return res, err
}

func (r *UserRepository) GetWithrawn(ctx context.Context, userId uint64) (res models.Money, err error) {
	return res, r.tx.GetContext(ctx, &res, `SELECT SUM("sum") as withdrawn FROM withdrawals WHERE user_id = $1`, userId)
}
