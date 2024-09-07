package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID UserID `json:"user_id,omitempty"`
}

type JWTManager struct {
	secret string
	Lifetime time.Duration
}

func NewJWTManager(secret string, lifetime time.Duration) *JWTManager {
	return &JWTManager{secret, lifetime}
}

func (jm *JWTManager) CreateTokenString(userID UserID) (string, error) {
	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jm.Lifetime)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(jm.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
} 
