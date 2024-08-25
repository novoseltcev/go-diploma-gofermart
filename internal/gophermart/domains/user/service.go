package user

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/novoseltcev/go-diploma-gofermart/internal/gophermart/models"
)


type UserStorager interface {
	IsLoginExists(ctx context.Context, login string) (bool, error)
	GetById(ctx context.Context, userId uint64) (*models.User, error)
	GetByLogin(ctx context.Context, login string) (*models.User, error)
	Create(ctx context.Context, login, password string) (uint64, error)
	GetWithrawn(ctx context.Context, userId uint64) (models.Money, error)
}

var (
	ErrAlreadyExists = errors.New("User already exists")
	ErrNotExists = errors.New("User not exists")
)

var _ UserStorager = &UserRepository{}


func Register(ctx context.Context, storager UserStorager, login, password string) (userId uint64, err error) {
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

	return user.Id, nil
}

func GetBalance(ctx context.Context, storager UserStorager, userId uint64) (*models.Balance, error) {
	user, err := storager.GetById(ctx, userId)
	if err != nil {
		return nil, err
	}
	
	if user == nil {
		return nil, ErrNotExists
	}

	withdrawn, err := storager.GetWithrawn(ctx, userId)
	if err != nil {
		return nil, err
	}

	return &models.Balance{
		Balance: user.Balance.Value,
		Withdrawn: withdrawn.Value,
	}, nil
}
