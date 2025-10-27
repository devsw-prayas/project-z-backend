package services

import (
	"database/sql"
	"errors"
	"log"
	"project-z-backend/database"
	"project-z-backend/models"
	"time"
)

func GetProblems() ([]models.Problem, error) {
	log.Println("GetProblems service called")

	rows, err := database.DB.Query(
		`SELECT id, title, description, difficulty, example_input, example_output, created_at 
		 FROM problems 
		 ORDER BY id ASC`,
	)
	if err != nil {
		log.Println("DB query error:", err)
		return nil, err
	}
	defer rows.Close()

	var problems []models.Problem

	for rows.Next() {
		var p models.Problem
		err := rows.Scan(&p.ID, &p.Title, &p.Description, &p.Difficulty, &p.ExampleInput, &p.ExampleOutput, &p.CreatedAt)
		if err != nil {
			log.Println("Row scan error:", err)
			return nil, err
		}

		// Fetch test cases for this problem
		testCases, err := getTestCasesByProblemID(p.ID)
		if err != nil {
			log.Printf("Error fetching test cases for problem %d: %v\n", p.ID, err)
		}
		p.TestCases = testCases

		problems = append(problems, p)
	}

	if err = rows.Err(); err != nil {
		log.Println("Row iteration error:", err)
		return nil, err
	}

	return problems, nil
}

func GetProblemByID(problemID int64) (*models.Problem, error) {
	log.Printf("GetProblemByID service called with ID: %d\n", problemID)

	var p models.Problem
	err := database.DB.QueryRow(
		`SELECT id, title, description, difficulty, example_input, example_output, created_at 
		 FROM problems 
		 WHERE id = $1`,
		problemID,
	).Scan(&p.ID, &p.Title, &p.Description, &p.Difficulty, &p.ExampleInput, &p.ExampleOutput, &p.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Problem with ID %d not found\n", problemID)
			return nil, nil
		}
		log.Println("DB query error:", err)
		return nil, err
	}

	// Fetch test cases for this problem
	testCases, err := getTestCasesByProblemID(problemID)
	if err != nil {
		log.Printf("Error fetching test cases for problem %d: %v\n", problemID, err)
	}
	p.TestCases = testCases

	return &p, nil
}

func CreateProblem(p models.Problem) (models.Problem, error) {
	log.Println("CreateProblem service called")

	// Validate required fields
	if p.Title == "" || p.Description == "" || p.Difficulty == "" {
		return models.Problem{}, errors.New("title, description, and difficulty are required")
	}

	// Insert problem into database
	err := database.DB.QueryRow(
		`INSERT INTO problems (title, description, difficulty, example_input, example_output, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, created_at`,
		p.Title, p.Description, p.Difficulty, p.ExampleInput, p.ExampleOutput, time.Now(),
	).Scan(&p.ID, &p.CreatedAt)

	if err != nil {
		log.Println("DB insert error:", err)
		return models.Problem{}, err
	}

	return p, nil
}

// Helper function to fetch test cases for a problem
func getTestCasesByProblemID(problemID int64) ([]models.TestCase, error) {
	log.Printf("getTestCasesByProblemID service called with problem ID: %d\n", problemID)

	rows, err := database.DB.Query(
		`SELECT id, problem_id, input, expected_output, is_sample, created_at 
		 FROM test_cases 
		 WHERE problem_id = $1 
		 ORDER BY id ASC`,
		problemID,
	)
	if err != nil {
		log.Printf("DB query error for test cases: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var testCases []models.TestCase

	for rows.Next() {
		var tc models.TestCase
		err := rows.Scan(&tc.ID, &tc.ProblemID, &tc.Input, &tc.ExpectedOutput, &tc.IsSample, &tc.CreatedAt)
		if err != nil {
			log.Printf("Row scan error for test case: %v\n", err)
			return nil, err
		}
		testCases = append(testCases, tc)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Row iteration error for test cases: %v\n", err)
		return nil, err
	}

	return testCases, nil
}
