package responses

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ValidationErr(c *gin.Context, err error) {
	_ = c.Error(err)
	c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"msg": err.Error()})
}
