package endpoints

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/novoseltcev/go-diploma-gofermart/gophermart/domains/orders"
	"github.com/novoseltcev/go-diploma-gofermart/gophermart/domains/orders/storager"
	r "github.com/novoseltcev/go-diploma-gofermart/gophermart/responses"
	"github.com/novoseltcev/go-diploma-gofermart/gophermart/utils"
	"github.com/novoseltcev/go-diploma-gofermart/shared"
)



func AddOrder(uowPool shared.UOWPool) gin.HandlerFunc {
	type reqBody struct {
		Number string `json:"number" binding:"required"`
	}

	return func (c *gin.Context) {
		userID := utils.GetUserID(c)

		var body reqBody
		if err := c.Bind(&body); err != nil {
			r.ErrInvalidRequest(c, err)
			return
		}

		uow := uowPool(c)
		defer uow.Close()

		if err := orders.AddOrderToUser(c, storager.New(uow), userID, body.Number); err != nil {
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
