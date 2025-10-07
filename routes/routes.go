package routes

import (
	"github.com/gin-gonic/gin"
)

func SetupAPIRoutes(router *gin.Engine) {
	api := router.Group("/api")
	setupUserRoutes(api)

}