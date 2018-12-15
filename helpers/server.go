package helpers

import "github.com/gin-gonic/gin"

func RespondWithError(code int, message string, c *gin.Context) {
	resp := map[string]string{
		"status":        "error",
		"error_message": message,
	}

	c.JSON(code, resp)
	c.Abort()
}
