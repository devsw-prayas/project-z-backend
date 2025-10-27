package routes

import (
	"project-z-backend/controllers"
	"project-z-backend/middleware"
	"project-z-backend/services"

	"github.com/gin-gonic/gin"
)

func setupHealthRoutes(api *gin.RouterGroup) {
	api.GET("/health", services.HealthCheck)
}

func setupUserRoutes(api *gin.RouterGroup) {

	api.POST("/user/register", controllers.Register)
	api.POST("/user/login", controllers.Login)
	api.GET("/user/me", middleware.AuthMiddleware(), controllers.UserInfo)
}
