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

		// Fetch statistics for this problem
		stats, err := getProblemStatistics(p.ID)
		if err != nil {
			log.Printf("Error fetching statistics for problem %d: %v\n", p.ID, err)
		}
		p.TotalSubmissions = stats.TotalSubmissions
		p.SuccessfulSubmissions = stats.SuccessfulSubmissions
		p.FailedSubmissions = stats.FailedSubmissions
		p.SuccessRate = stats.SuccessRate
		p.FailureRate = stats.FailureRate

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

	// Fetch statistics for this problem
	stats, err := getProblemStatistics(problemID)
	if err != nil {
		log.Printf("Error fetching statistics for problem %d: %v\n", problemID, err)
	}
	p.TotalSubmissions = stats.TotalSubmissions
	p.SuccessfulSubmissions = stats.SuccessfulSubmissions
	p.FailedSubmissions = stats.FailedSubmissions
	p.SuccessRate = stats.SuccessRate
	p.FailureRate = stats.FailureRate

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

// ProblemStatistics holds submission statistics for a problem
type ProblemStatistics struct {
	TotalSubmissions      int64
	SuccessfulSubmissions int64
	FailedSubmissions     int64
	SuccessRate           float64
	FailureRate           float64
}

// Helper function to calculate problem statistics
func getProblemStatistics(problemID int64) (ProblemStatistics, error) {
	log.Printf("getProblemStatistics service called with problem ID: %d\n", problemID)

	stats := ProblemStatistics{}

	// Count total submissions
	err := database.DB.QueryRow(
		`SELECT COUNT(*) FROM submissions WHERE problem_id = $1`,
		problemID,
	).Scan(&stats.TotalSubmissions)
	if err != nil {
		log.Printf("Error counting total submissions: %v\n", err)
		return stats, err
	}

	// Count successful submissions
	err = database.DB.QueryRow(
		`SELECT COUNT(*) FROM submissions WHERE problem_id = $1 AND status = 'Accepted'`,
		problemID,
	).Scan(&stats.SuccessfulSubmissions)
	if err != nil {
		log.Printf("Error counting successful submissions: %v\n", err)
		return stats, err
	}

	// Count failed submissions
	stats.FailedSubmissions = stats.TotalSubmissions - stats.SuccessfulSubmissions

	// Calculate rates
	if stats.TotalSubmissions > 0 {
		stats.SuccessRate = (float64(stats.SuccessfulSubmissions) / float64(stats.TotalSubmissions)) * 100
		stats.FailureRate = (float64(stats.FailedSubmissions) / float64(stats.TotalSubmissions)) * 100
	}

	return stats, nil
}

// CreateSubmission creates a new submission for a user solving a problem
func CreateSubmission(submission models.Submission) (models.Submission, error) {
	log.Println("CreateSubmission service called")

	// Validate required fields
	if submission.UserID == 0 || submission.ProblemID == 0 || submission.Language == "" || submission.SourceCode == "" {
		return models.Submission{}, errors.New("user_id, problem_id, language, and source_code are required")
	}

	// Insert submission into database
	err := database.DB.QueryRow(
		`INSERT INTO submissions (user_id, problem_id, language, source_code, status, submitted_at)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, submitted_at`,
		submission.UserID, submission.ProblemID, submission.Language, submission.SourceCode, "Queued", time.Now(),
	).Scan(&submission.ID, &submission.SubmittedAt)

	if err != nil {
		log.Printf("DB insert error: %v\n", err)
		return models.Submission{}, err
	}

	submission.Status = "Queued"
	return submission, nil
}

// UpdateSubmissionStatus updates a submission's status with result details
func UpdateSubmissionStatus(submissionID int64, status string, resultMessage string, executionTime float64, memoryUsed float64) error {
	log.Printf("UpdateSubmissionStatus service called with ID: %d, Status: %s\n", submissionID, status)

	_, err := database.DB.Exec(
		`UPDATE submissions 
		 SET status = $1, result_message = $2, execution_time = $3, memory_used = $4 
		 WHERE id = $5`,
		status, resultMessage, executionTime, memoryUsed, submissionID,
	)

	if err != nil {
		log.Printf("DB update error: %v\n", err)
		return err
	}

	return nil
}

// GetSubmissionByID fetches a specific submission
func GetSubmissionByID(submissionID int64) (*models.Submission, error) {
	log.Printf("GetSubmissionByID service called with ID: %d\n", submissionID)

	var s models.Submission
	err := database.DB.QueryRow(
		`SELECT id, user_id, problem_id, language, source_code, status, result_message, execution_time, memory_used, submitted_at 
		 FROM submissions 
		 WHERE id = $1`,
		submissionID,
	).Scan(&s.ID, &s.UserID, &s.ProblemID, &s.Language, &s.SourceCode, &s.Status, &s.ResultMessage, &s.ExecutionTime, &s.MemoryUsed, &s.SubmittedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Submission with ID %d not found\n", submissionID)
			return nil, nil
		}
		log.Printf("DB query error: %v\n", err)
		return nil, err
	}

	return &s, nil
}
