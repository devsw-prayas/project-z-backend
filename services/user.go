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
	if u.Name == "" || u.Email == "" || u.Password == "" {
		return models.User{}, errors.New("username, email, and password are required")
	}

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

func UserInfo(userID int64) (models.User, error) {
	var u models.User
	err := database.DB.QueryRow(
		`SELECT id, username, email, created_at FROM users WHERE id = $1`,
		userID,
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

// GetUserWithStats fetches user information including submission stats and solved problems
func GetUserWithStats(userID int64) (*models.User, error) {
	log.Printf("GetUserWithStats service called with ID: %d\n", userID)

	var u models.User
	err := database.DB.QueryRow(
		`SELECT id, username, email, created_at FROM users WHERE id = $1`,
		userID,
	).Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt)

	if err != nil {
		log.Printf("Error fetching user: %v\n", err)
		return nil, err
	}

	// Get total submissions count
	err = database.DB.QueryRow(
		`SELECT COUNT(*) FROM submissions WHERE user_id = $1`,
		userID,
	).Scan(&u.TotalSubmissions)
	if err != nil {
		log.Printf("Error counting submissions: %v\n", err)
	}

	// Get solved problems count
	err = database.DB.QueryRow(
		`SELECT COUNT(DISTINCT problem_id) FROM submissions WHERE user_id = $1 AND status = 'Accepted'`,
		userID,
	).Scan(&u.SolvedProblems)
	if err != nil {
		log.Printf("Error counting solved problems: %v\n", err)
	}

	// Get failed problems count
	err = database.DB.QueryRow(
		`SELECT COUNT(DISTINCT problem_id) FROM submissions WHERE user_id = $1 AND status != 'Accepted'`,
		userID,
	).Scan(&u.FailedProblems)
	if err != nil {
		log.Printf("Error counting failed problems: %v\n", err)
	}

	// Fetch user's submissions
	submissions, err := GetUserSubmissions(userID)
	if err != nil {
		log.Printf("Error fetching submissions: %v\n", err)
	}
	u.Submissions = submissions

	u.Password = ""
	return &u, nil
}

// GetUserSubmissions fetches all submissions for a user
func GetUserSubmissions(userID int64) ([]models.Submission, error) {
	log.Printf("GetUserSubmissions service called with user ID: %d\n", userID)

	rows, err := database.DB.Query(
		`SELECT id, user_id, problem_id, language, source_code, status, result_message, execution_time, memory_used, submitted_at 
		 FROM submissions 
		 WHERE user_id = $1 
		 ORDER BY submitted_at DESC`,
		userID,
	)
	if err != nil {
		log.Printf("DB query error: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var submissions []models.Submission

	for rows.Next() {
		var s models.Submission
		err := rows.Scan(&s.ID, &s.UserID, &s.ProblemID, &s.Language, &s.SourceCode, &s.Status, &s.ResultMessage, &s.ExecutionTime, &s.MemoryUsed, &s.SubmittedAt)
		if err != nil {
			log.Printf("Row scan error: %v\n", err)
			return nil, err
		}
		submissions = append(submissions, s)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Row iteration error: %v\n", err)
		return nil, err
	}

	return submissions, nil
}
