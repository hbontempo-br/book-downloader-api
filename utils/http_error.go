package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func DefaultErrorMessage(c *gin.Context, httpStatus int, details interface{}) {
	resp := map[string]interface{}{
		"code":    httpStatus,
		"message": http.StatusText(httpStatus),
		"details": details,
	}
	c.JSON(httpStatus, resp)
}
