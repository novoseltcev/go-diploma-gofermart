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
		userId := utils.GetUserId(c)

		uow := uowPool(c)
		defer uow.Close()

		balance, err := balance.GetBalance(c, storager.New(uow), userId)

		if err != nil {
			r.InternalErr(c, err)
			return
		}
		c.JSON(http.StatusOK, balance)
	}
}
