package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

var DB *sql.DB

func InitDB(databaseUrL string) {
	if databaseUrL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	var err error
	DB, err = sql.Open("postgres", databaseUrL)
	if err != nil {
		log.Fatal("Failed to open DB:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Cannot connect to DB:", err)
	}

	log.Println("Connected to DB")
}
