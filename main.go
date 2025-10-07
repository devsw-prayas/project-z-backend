package main

import (
	"log"
	"project-z-backend/config"
	"project-z-backend/database"
	"project-z-backend/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	cfg := config.LoadConfig()

	database.InitDB(cfg)
	defer database.DB.Close()

	router := gin.Default()
	routes.SetupAPIRoutes(router)

	log.Printf("🚀 Server is running on port %s", cfg.Port)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
