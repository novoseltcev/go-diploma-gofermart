package endpoints

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/novoseltcev/go-diploma-gofermart/gophermart/domains/balance"
	"github.com/novoseltcev/go-diploma-gofermart/gophermart/domains/balance/storager"
	r "github.com/novoseltcev/go-diploma-gofermart/gophermart/responses"
	"github.com/novoseltcev/go-diploma-gofermart/gophermart/utils"
	"github.com/novoseltcev/go-diploma-gofermart/shared"
)


func Withdraw(uowPool shared.UOWPool) gin.HandlerFunc {
	type reqBody struct {
		Sum float32  `json:"sum" binding:"required"`
		Order string `json:"order" binding:"required"`
	}

	return func (c *gin.Context) {
		userID := utils.GetUserID(c)
		
		var body reqBody
		if err := c.BindJSON(&body); err != nil {
			r.ErrInvalidRequest(c, err)
			return
		}

		uow := uowPool(c)
		defer uow.Close()

		err := balance.Withdrawn(c, storager.New(uow), userID, body.Sum, body.Order)
		if err != nil {
			if errors.Is(err, balance.ErrLunhNumberValidation) {
				r.ErrValidation(c, err)
			} else if errors.Is(err, balance.ErrNotEnought) {
				_ = c.Error(err)
				c.AbortWithStatusJSON(http.StatusPaymentRequired, gin.H{"msg": err.Error()})
			}
			return
		}
		_ = uow.Apply()
		c.JSON(http.StatusOK, nil)
	}
}
