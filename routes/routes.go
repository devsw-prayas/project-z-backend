package routes

import (
	"project-z-backend/handlers"

	"github.com/gin-gonic/gin"
)

func SetupAPIRoutes(router *gin.Engine) {
	router.GET("/", handlers.WelcomeHandler)

	api := router.Group("/api")
	api.GET("/health", handlers.HealthHandler)
	setupUserRoutes(api)
	setupProblemRoutes(api)
}
