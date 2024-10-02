package responses

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrLogic(c *gin.Context, err error) {
	_ = c.Error(err)
	c.AbortWithStatusJSON(http.StatusConflict, gin.H{"msg": err.Error()})
}
