package endpoints

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/novoseltcev/go-diploma-gofermart/internal/gophermart/domains/balance"
	"github.com/novoseltcev/go-diploma-gofermart/internal/gophermart/domains/balance/storager"
	r "github.com/novoseltcev/go-diploma-gofermart/internal/gophermart/responses"
	"github.com/novoseltcev/go-diploma-gofermart/internal/gophermart/utils"
	"github.com/novoseltcev/go-diploma-gofermart/internal/shared"
)


func Withdraw(uowPool shared.UOWPool) gin.HandlerFunc {
	type reqBody struct {
		Sum uint64  `json:"sum" binding:"required"`
		Order string `json:"order" binding:"required"`
	}

	return func (c *gin.Context) {
		userId := utils.GetUserId(c)
		
		var body reqBody
		if err := c.BindJSON(&body); err != nil {
			r.InvalidRequestErr(c, err)
			return
		}

		uow := uowPool(c)
		defer uow.Close()

		err := balance.Withdrawn(c, storager.New(uow), userId, body.Sum, body.Order)
		if err != nil {
			if errors.Is(err, balance.LunhNumberValidationErr) {
				r.ValidationErr(c, err)
			} else if errors.Is(err, balance.NotEnoughtErr) {
				_ = c.Error(err)
				c.AbortWithStatusJSON(http.StatusPaymentRequired, gin.H{"msg": err.Error()})
			}
			return
		}
		_ = uow.Apply()
		c.JSON(http.StatusOK, nil)
	}
}
