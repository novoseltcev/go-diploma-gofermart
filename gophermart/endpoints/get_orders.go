package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/novoseltcev/go-diploma-gofermart/gophermart/domains/orders"
	"github.com/novoseltcev/go-diploma-gofermart/gophermart/domains/orders/storager"
	r "github.com/novoseltcev/go-diploma-gofermart/gophermart/responses"
	"github.com/novoseltcev/go-diploma-gofermart/gophermart/utils"
	"github.com/novoseltcev/go-diploma-gofermart/shared"
)


func GetOrders(uowPool shared.UOWPool) gin.HandlerFunc {
	return func (c *gin.Context) {
		userID := utils.GetUserID(c)
		uow := uowPool(c)
		defer uow.Close()

		orders, err := orders.GetUserOrders(c, storager.New(uow), userID)
		if err != nil {
			r.ErrInternal(c, err)
			return
		}

		if len(orders) == 0 {
			c.JSON(http.StatusNoContent, nil)
		} else {
			c.JSON(http.StatusOK, orders)
		}
	}
}
