package handlers

import (
	"github.com/gin-gonic/gin"
)

func HandleError(c *gin.Context, statusCode int, errorMsg string) {
	msg := errorMsg
	if msg == "" {
		msg = "Internal Server Error"
	}
	c.JSON(statusCode, gin.H{"error": msg})

}

func HandleSuccess(c *gin.Context, statusCode int, responseData interface{}) {
	c.JSON(statusCode, responseData)
}
