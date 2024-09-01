package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	jwt.RegisteredClaims
	UserId UserId `json:"user_id,omitempty"`
}

type JWTManager struct {
	secret string
	Lifetime time.Duration
}

func NewJWTManager(secret string, lifetime time.Duration) *JWTManager {
	return &JWTManager{secret, lifetime}
}

func (jm *JWTManager) CreateTokenString(userId UserId) (string, error) {
	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jm.Lifetime)),
		},
		UserId: userId,
	})

	tokenString, err := token.SignedString([]byte(jm.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
} 
