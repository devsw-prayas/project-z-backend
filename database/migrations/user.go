package migrations

import (
	"log"
	"os"
	"project-z-backend/database"
)

func UsersMigration() {
	if database.DB == nil {
		log.Fatal("Database not initialized")
	}

	bytefile, err := os.ReadFile("database/users.sql")
	if err != nil {
		log.Fatal("Failed to read users.sql file:", err)
	}

	database.DB.Exec(string(bytefile))
	if err != nil {
		log.Fatal("Failed to create users table:", err)
	}

	log.Println("Database migrations ran successfully")
}
