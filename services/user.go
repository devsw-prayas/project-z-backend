package services

import (
	"database/sql"
	_ "database/sql"
	"errors"
	"time"

	"project-z-backend/database"
	"project-z-backend/middleware"
	"project-z-backend/models"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func Register(u models.User) (models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}

	err = database.DB.QueryRow(
		`INSERT INTO users (username, email, password_hash, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, username, email, created_at`,
		u.Name, u.Email, string(hashedPassword), time.Now(),
	).Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt)

	if err != nil {
		return models.User{}, err
	}

	return u, nil
}

func UserInfo(u models.User) (models.User, error) {

	err := database.DB.QueryRow(
		`SELECT id, username, email, created_at FROM users WHERE username = $1`,
		u.ID,
	).Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt)
	if err != nil {
		return models.User{}, err
	}

	return u, nil
}

func Login(u models.User) (string, error) {
	var passwordHash string
	err := database.DB.QueryRow(
		`SELECT id, username, email, password_hash FROM users WHERE username = $1`,
		u.Name,
	).Scan(&u.ID, &u.Name, &u.Email, &passwordHash)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("Invalid username or password")
		}
		return "", errors.New("Database error")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(u.Password)); err != nil {
		return "", errors.New("Invalid username or password")
	}

	tokenString, err := middleware.CreateJWT(u.ID, u.Email)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
