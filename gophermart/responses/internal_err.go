package responses

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func InternalErr(c *gin.Context, err error) {
	_ = c.Error(err)
	c.AbortWithStatus(http.StatusInternalServerError)
}
