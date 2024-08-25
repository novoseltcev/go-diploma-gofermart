package endpoints

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/novoseltcev/go-diploma-gofermart/internal/gophermart/auth"
	"github.com/novoseltcev/go-diploma-gofermart/internal/gophermart/domains/users"
	"github.com/novoseltcev/go-diploma-gofermart/internal/gophermart/domains/users/storager"
	r "github.com/novoseltcev/go-diploma-gofermart/internal/gophermart/responses"
	"github.com/novoseltcev/go-diploma-gofermart/internal/shared"
)


func Login(uowPool shared.UOWPool, jwtManager *auth.JWTManager) gin.HandlerFunc {
	type reqBody struct {
		Login string `json:"login" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	return func (c *gin.Context) {
		var body reqBody
		if err := c.Bind(&body); err != nil {
			r.InvalidRequestErr(c, err)
			return
		}

		uow := uowPool(c)
		defer uow.Close()

		userId, err := users.Login(c, storager.New(uow), body.Login, body.Password)
		if err != nil {
			if errors.Is(err, users.ErrNotExists) {
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
