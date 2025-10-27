package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Submit(c *gin.Context) {
	// TODO: Handle new submission
	c.JSON(http.StatusOK, gin.H{"message": "submission received"})
}

func GetSubmissionID(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"submission_id": id, "status": "queued"})
}

func Executor(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"executors": []string{"cpp", "python", "go"}})
}
