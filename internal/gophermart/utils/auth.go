package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/novoseltcev/go-diploma-gofermart/internal/gophermart/auth"
)


func GetUserId(c *gin.Context) auth.UserId {
	authIdentity, ok := c.Get(auth.IdentityKey)
	if !ok {
		panic(errors.New("auth identity not passed to router"))
	}

	userId, ok := authIdentity.(auth.UserId)
	if !ok {
		panic(errors.New("auth identity not passed to router"))
	}
	return userId
}
