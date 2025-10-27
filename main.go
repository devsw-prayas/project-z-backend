package main

import (
	"log"
	"project-z-backend/config"
	"project-z-backend/database"
	"project-z-backend/database/migrations"
	"project-z-backend/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	cfg := config.LoadConfig()

	database.InitDB(cfg.DB_URL)
	defer database.DB.Close()

	migrations.SetupMigration()

	router := gin.Default()
	routes.SetupAPIRoutes(router)

	log.Printf("Server is running on port %s", cfg.PORT)

	if err := router.Run(":" + cfg.PORT); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
