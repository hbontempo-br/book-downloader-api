package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func DefaultErrorMessage(c *gin.Context, httpStatus int, details interface{}) {
	resp := map[string]interface{}{
		"code":    httpStatus,
		"message": http.StatusText(httpStatus),
		"details": details,
	}
	c.JSON(httpStatus, resp)
}
