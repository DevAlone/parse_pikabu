package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func AnswerError(c *gin.Context, statusCode int, message string) {
	switch statusCode {
	case 404:
		message = "not found. " + message
	case 500:
		message = "internal server error. " + message
	}

	c.JSON(statusCode, map[string]string{
		"status":      "error",
		"status_code": fmt.Sprint(statusCode),
		"message":     message,
	})
}
