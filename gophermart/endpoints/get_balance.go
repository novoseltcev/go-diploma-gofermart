package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/novoseltcev/go-diploma-gofermart/gophermart/domains/balance"
	"github.com/novoseltcev/go-diploma-gofermart/gophermart/domains/balance/storager"
	r "github.com/novoseltcev/go-diploma-gofermart/gophermart/responses"
	"github.com/novoseltcev/go-diploma-gofermart/gophermart/utils"
	"github.com/novoseltcev/go-diploma-gofermart/shared"
)


func GetBalance(uowPool shared.UOWPool) gin.HandlerFunc {
	return func (c *gin.Context) {
		userID := utils.GetUserID(c)

		uow := uowPool(c)
		defer uow.Close()

		balance, err := balance.GetBalance(c, storager.New(uow), userID)

		if err != nil {
			r.ErrInternal(c, err)
			return
		}

		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, balance)
	}
}
