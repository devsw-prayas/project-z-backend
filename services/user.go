package services

import (
	_ "database/sql"
	"errors"
	"log"
	"os"
	"project-z-backend/database"
	"project-z-backend/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func Register(u models.User) (models.User, error) {
	log.Println("Register handler called")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}

	// log.Printf("User data received: %+v\n", u)

	err = database.DB.QueryRow(
		`INSERT INTO users (username, email, password_hash, created_at)
     	VALUES ($1, $2, $3, $4)
     	RETURNING id, created_at`,
		u.Name, u.Email, string(hashedPassword), time.Now(),
	).Scan(&u.ID, &u.CreatedAt)

	if err != nil {
		log.Println("DB insert error:", err)
		return models.User{}, err
	}

	u.Password = ""

	return u, nil
}

func UserInfo(u models.User) (models.User, error) {

	log.Println("UserInfo handler called")

	err := database.DB.QueryRow(
		`SELECT id, username, email, created_at FROM users WHERE username = $1`,
		u.Name,
	).Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt)
	if err != nil {
		return models.User{}, err
	}
	u.Password = ""

	return u, nil
}

func Login(u models.User) (string, error) {

	log.Println("Login handler called")

	var passwordHash string
	err := database.DB.QueryRow(
		`SELECT id, username, email, password_hash FROM users WHERE username = $1`,
		u.Name,
	).Scan(&u.ID, &u.Name, &u.Email, &passwordHash)
	if err != nil {
		return "", errors.New("Invalid username or password")
	}

	if bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(u.Password)) != nil {
		return "", errors.New("Invalid username or password")
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": u.ID,
		"email":   u.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	secret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", errors.New("Failed to generate token")
	}

	return tokenString, nil
}
