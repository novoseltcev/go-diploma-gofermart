package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/novoseltcev/go-diploma-gofermart/gophermart/auth"
)


func GetUserID(c *gin.Context) auth.UserID {
	authIdentity, ok := c.Get(auth.IdentityKey)
	if !ok {
		panic(errors.New("auth identity not passed to router"))
	}

	userID, ok := authIdentity.(auth.UserID)
	if !ok {
		panic(errors.New("auth identity not passed to router"))
	}
	return userID
}
