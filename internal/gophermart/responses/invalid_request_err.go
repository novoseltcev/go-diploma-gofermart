package responses

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func InvalidRequestErr(c *gin.Context, err error) {
	_ = c.Error(err)
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
}
