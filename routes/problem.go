package routes

import (
	"project-z-backend/controllers"

	"github.com/gin-gonic/gin"
)

func setupProblemRoutes(api *gin.RouterGroup) {
	api.GET("/problems", controllers.GetProblems)
	api.POST("/problems", controllers.CreateProblem)
	api.GET("/problems/:problem_id", controllers.GetProblemByID)
}
