package migrations
import (
	"log"
	"project-z-backend/database"
)


func UsersMigration() {
	if database.DB == nil {
		log.Fatal("Database not initialized")
	}

	_, err := database.DB.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id BIGSERIAL PRIMARY KEY,
		username VARCHAR(50) UNIQUE NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`)
	if err != nil {
		log.Fatal("Failed to create users table:", err)
	}

	log.Println("✅ Database migrations ran successfully")
}