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




func Register(uowPool shared.UOWPool, jwtManager *auth.JWTManager) gin.HandlerFunc {
	return func (c *gin.Context) {
		var authData schemas.AuthData
		if err := c.Bind(&authData); err != nil {
			r.InvalidRequestErr(c, err)
			return
		}

		uow := uowPool(c)
		defer uow.Close()

		userId, err := domain.Register(c, domain.NewStorager(uow), authData.Login, authData.Password)
		if err != nil {
			if errors.Is(err, domain.ErrAlreadyExists) {
				r.LogicErr(c, err)
			} else {
				r.InternalErr(c, err)
			}
			return
		}

		_ = uow.Apply()
	
		token, err := jwtManager.CreateTokenString(userId)
		if err != nil {
			r.InternalErr(c, err)
		}

		c.JSON(http.StatusOK, gin.H{"access": token})
	}
}
