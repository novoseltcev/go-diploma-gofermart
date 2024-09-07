package users

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/novoseltcev/go-diploma-gofermart/gophermart/models"
)


type UserStorager interface {
	IsLoginExists(ctx context.Context, login string) (bool, error)
	GetByLogin(ctx context.Context, login string) (*models.User, error)
	Create(ctx context.Context, login, password string) (uint64, error)
}


var (
	ErrAlreadyExists = errors.New("user already exists")
	ErrNotExists = errors.New("user not exists")
)


func Register(ctx context.Context, storager UserStorager, login, password string) (userID uint64, err error) {
	exists, err := storager.IsLoginExists(ctx, login)
	if err != nil {
		return 0, err
	}

	if exists {
		return 0, ErrAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	return storager.Create(ctx, login, string(hashedPassword))
}

func Login(ctx context.Context, storager UserStorager, login, password string) (uint64, error) {
	user, err := storager.GetByLogin(ctx, login)
	if err != nil {
		return 0, err
	}

	if user == nil {
		return 0, ErrNotExists
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
		return 0, ErrNotExists
	}

	return user.ID, nil
}
