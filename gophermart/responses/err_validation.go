package responses

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrValidation(c *gin.Context, err error) {
	_ = c.Error(err)
	c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"msg": err.Error()})
}
