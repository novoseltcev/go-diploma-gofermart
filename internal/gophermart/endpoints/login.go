package endpoints

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/novoseltcev/go-diploma-gofermart/internal/gophermart/auth"
	domain "github.com/novoseltcev/go-diploma-gofermart/internal/gophermart/domains/user"
	r "github.com/novoseltcev/go-diploma-gofermart/internal/gophermart/responses"
	"github.com/novoseltcev/go-diploma-gofermart/internal/gophermart/schemas"
	"github.com/novoseltcev/go-diploma-gofermart/internal/shared"
)


func Login(uowPool shared.UOWPool, jwtManager *auth.JWTManager) gin.HandlerFunc {
	return func (c *gin.Context) {
		var scheme schemas.AuthData
		if err := c.Bind(&scheme); err != nil {
			r.InvalidRequestErr(c, err)
			return
		}

		uow := uowPool(c)
		defer uow.Close()

		userId, err := domain.Login(c, domain.NewStorager(uow), scheme.Login, scheme.Password)
		if err != nil {
			if errors.Is(err, domain.ErrNotExists) {
				r.UnauthorizedErr(c, err)
			} else {
				r.InternalErr(c, err)
			}
			return
		}

		token, err := jwtManager.CreateTokenString(userId)
		if err != nil {
			r.InternalErr(c, err)
		}

		c.JSON(http.StatusOK, gin.H{"accessToken": token})
	}
}
