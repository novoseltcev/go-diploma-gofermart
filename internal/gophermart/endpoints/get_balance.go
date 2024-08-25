package endpoints

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	domain "github.com/novoseltcev/go-diploma-gofermart/internal/gophermart/domains/user"
	r "github.com/novoseltcev/go-diploma-gofermart/internal/gophermart/responses"
	"github.com/novoseltcev/go-diploma-gofermart/internal/gophermart/utils"
	"github.com/novoseltcev/go-diploma-gofermart/internal/shared"
)


func GetBalance(uowPool shared.UOWPool) gin.HandlerFunc {
	return func (c *gin.Context) {
		userId := utils.GetUserId(c)

		uow := uowPool(c)
		defer uow.Close()

		balance, err := domain.GetBalance(c, domain.NewStorager(uow), userId)

		if err != nil {
			if errors.Is(err, domain.ErrNotExists) {
				r.UnauthorizedErr(c, err)
			} else {
				r.InternalErr(c, err)
			}
			return
		}
		c.JSON(http.StatusOK, balance)
	}
}
