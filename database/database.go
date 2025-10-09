package database

import (
	"database/sql"
	"log"
	"project-z-backend/config"


	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB(config *config.Config) {
	if config.DB_URL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	var err error
	DB, err = sql.Open("postgres", config.DB_URL)
	if err != nil {
		log.Fatal("Failed to open DB:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Cannot connect to DB:", err)
	}

	log.Println("✅ Connected to PostgreSQL")
}
