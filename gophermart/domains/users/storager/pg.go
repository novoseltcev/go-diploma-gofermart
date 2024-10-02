package storager

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"github.com/novoseltcev/go-diploma-gofermart/gophermart/domains/users"
	"github.com/novoseltcev/go-diploma-gofermart/gophermart/models"
	"github.com/novoseltcev/go-diploma-gofermart/shared"
)

type repository struct {
	tx *sqlx.Tx
} 

func New(uow *shared.UnitOfWork) users.UserStorager {
	return &repository{tx: uow.Tx}
}

func (r *repository) IsLoginExists(ctx context.Context, login string) (result bool, err error) {
	err = r.tx.GetContext(ctx, &result, "SELECT COUNT(login) > 0 FROM users WHERE login = $1", login)
	return result, err
}

func (r *repository) GetByLogin(ctx context.Context, login string) (result *models.User, err error) {
	err = r.tx.GetContext(ctx, result, "SELECT id, login, balance, hashed_password FROM users WHERE login = $1", login)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return result, err
}

func (r *repository) Create(ctx context.Context, login, password string) (models.UserID, error) {
	var res uint64
	err := r.tx.GetContext(ctx, &res, "INSERT INTO users (login, hashed_password) VALUES ($1, $2) RETURNING id", login, password)
	if err != nil {
		return 0, err
	}
	return res, err
}
