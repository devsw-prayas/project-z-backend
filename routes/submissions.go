package routes

import (
	"project-z-backend/services"

	"github.com/gin-gonic/gin"
)

func setupSubmissionRoutes(api *gin.RouterGroup) {
	api.POST("/submissions", services.Submit)
	api.GET("/submissions/:id", services.GetSubmissionID)
	api.GET("/executors", services.Executor)
}
