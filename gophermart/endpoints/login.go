package endpoints

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/novoseltcev/go-diploma-gofermart/gophermart/auth"
	"github.com/novoseltcev/go-diploma-gofermart/gophermart/domains/users"
	"github.com/novoseltcev/go-diploma-gofermart/gophermart/domains/users/storager"
	r "github.com/novoseltcev/go-diploma-gofermart/gophermart/responses"
	"github.com/novoseltcev/go-diploma-gofermart/shared"
)


func Login(uowPool shared.UOWPool, jwtManager *auth.JWTManager) gin.HandlerFunc {
	type reqBody struct {
		Login string `json:"login" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	return func (c *gin.Context) {
		var body reqBody
		if err := c.Bind(&body); err != nil {
			r.ErrInvalidRequest(c, err)
			return
		}

		uow := uowPool(c)
		defer uow.Close()

		userID, err := users.Login(c, storager.New(uow), body.Login, body.Password)
		if err != nil {
			if errors.Is(err, users.ErrNotExists) {
				r.ErrUnauthorized(c, err)
			} else {
				r.ErrInternal(c, err)
			}
			return
		}

		token, err := jwtManager.CreateTokenString(userID)
		if err != nil {
			r.ErrInternal(c, err)
		}

		c.JSON(http.StatusOK, gin.H{"accessToken": token})
	}
}
