package endpoints

import (
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/novoseltcev/go-diploma-gofermart/gophermart/domains/orders"
	"github.com/novoseltcev/go-diploma-gofermart/gophermart/domains/orders/storager"
	r "github.com/novoseltcev/go-diploma-gofermart/gophermart/responses"
	"github.com/novoseltcev/go-diploma-gofermart/gophermart/utils"
	"github.com/novoseltcev/go-diploma-gofermart/shared"
)



func AddOrder(uowPool shared.UOWPool) gin.HandlerFunc {
	return func (c *gin.Context) {
		userID := utils.GetUserID(c)

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			r.ErrInvalidRequest(c, err)
			return
		}
		number := string(body)

		uow := uowPool(c)
		defer uow.Close()

		if err := orders.AddOrderToUser(c, storager.New(uow), userID, number); err != nil {
			if errors.Is(err, orders.ErrOrderLoaded) {
				c.JSON(http.StatusOK, nil)
			} else if errors.Is(err, orders.ErrLunhNumberValidation) {
				r.ErrValidation(c, err)
			} else if errors.Is(err, orders.ErrOrderNotMeLoaded) {
				r.ErrLogic(c, err)
			} else {
				r.ErrInternal(c, err)
			}
			return
		}
		_ = uow.Apply()

		c.JSON(http.StatusAccepted, nil)
	}
}
