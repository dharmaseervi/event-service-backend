package utils

import "github.com/gin-gonic/gin"

// RespondWithJSON sends a JSON response
func RespondWithJSON(c *gin.Context, statusCode int, payload interface{}) {
	c.JSON(statusCode, gin.H{
		"status": "success",
		"data":   payload,
	})
}

// RespondWithError sends an error response
func RespondWithError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{
		"status":  "error",
		"message": message,
	})
}
