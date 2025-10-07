package routes

import (
	"project-z-backend/services"

	"github.com/gin-gonic/gin"
)

func setupUserRoutes(api *gin.RouterGroup) {
	
	api.POST("/user/register", services.Register)
	api.POST("/user/login", services.Login)
	api.GET("/users/me", services.UserInfo)

}