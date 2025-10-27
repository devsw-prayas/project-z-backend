package migrations

import (
	"log"
	"os"
	"project-z-backend/database"
)

func ProblemsMigration() {
	if database.DB == nil {
		log.Fatal("Database not initialized")
	}

	byteFile, err := os.ReadFile("database/problems.sql")
	if err != nil {
		log.Fatal("Failed to read problems.sql", err)
	}

	_, err = database.DB.Exec(string(byteFile))
	if err != nil {
		log.Fatal("Failed to create or verify problems table", err)
	}

	log.Println("Database migration for Problems Ran successfully")

}
