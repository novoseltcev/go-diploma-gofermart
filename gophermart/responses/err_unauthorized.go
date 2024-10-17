package responses

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrUnauthorized(c *gin.Context, err error) {
	_ = c.Error(err)
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": err.Error()})
}
