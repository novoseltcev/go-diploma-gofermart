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


func Register(uowPool shared.UOWPool, jwtManager *auth.JWTManager) gin.HandlerFunc {
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

		userID, err := users.Register(c, storager.New(uow), body.Login, body.Password)
		if err != nil {
			if errors.Is(err, users.ErrAlreadyExists) {
				r.ErrLogic(c, err)
			} else {
				r.ErrInternal(c, err)
			}
			return
		}

		_ = uow.Apply()
	
		token, err := jwtManager.CreateTokenString(userID)
		if err != nil {
			r.ErrInternal(c, err)
		}

		c.Header("Authorization", token)
		c.JSON(http.StatusOK, gin.H{})
	}
}
