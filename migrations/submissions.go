package migrations

import (
	"log"
	"os"
	"project-z-backend/database"
)

func SubmissionsMigration() {
	if database.DB == nil {
		log.Fatal("Database not initialized")
	}

	byteFile, err := os.ReadFile("database/submissions.sql")
	if err != nil {
		log.Fatal("Failed to read submissions.sql", err)
	}

	_, err = database.DB.Exec(string(byteFile))
	if err != nil {
		log.Fatal("Failed to create or verify submissions table", err)
	}

	log.Println("Database migration for Submissions Ran successfully")
}
