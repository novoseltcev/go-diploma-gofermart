package auth

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"

	r "github.com/novoseltcev/go-diploma-gofermart/gophermart/responses"
)

func JWTMiddleware(manager *JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			r.ErrUnauthorized(c, errors.New(""))
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims,
			func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("token signed method is invalid")
				}
				return []byte(manager.secret), nil
			},
		)

		if err != nil {
			var validationErr jwt.ValidationError
			if errors.As(err, &validationErr) {
				r.ErrUnauthorized(c, validationErr.Unwrap())
			} else {
				r.ErrUnauthorized(c, err)
			}
			return
		}

		if !token.Valid {
			r.ErrUnauthorized(c, jwt.ErrTokenNotValidYet)
			return
		}

		if claims.UserID == 0 {
			r.ErrUnauthorized(c, jwt.ErrTokenInvalidClaims)
			return 
		}

		c.Set(IdentityKey, claims.UserID)
		c.Next()
	}
}
