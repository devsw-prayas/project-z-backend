package models

type User struct {
	ID                 int64        `json:"id"`
	Name               string       `json:"name"`
	Email              string       `json:"email"`
	Password           string       `json:"password"`
	CreatedAt          string       `json:"created_at"`
	TotalSubmissions   int64        `json:"total_submissions"`
	SolvedProblems     int64        `json:"solved_problems"`
	FailedProblems     int64        `json:"failed_problems"`
	Submissions        []Submission `json:"submissions"`
	SolvedProblemsList []Problem    `json:"solved_problems_list"`
}

type Submission struct {
	ID            int64   `json:"id"`
	UserID        int64   `json:"user_id"`
	ProblemID     int64   `json:"problem_id"`
	Language      string  `json:"language"`
	SourceCode    string  `json:"source_code"`
	Status        string  `json:"status"` // "Accepted", "Wrong Answer", "Runtime Error", "Time Limit Exceeded", "Compilation Error", "Queued"
	ResultMessage string  `json:"result_message"`
	ExecutionTime float64 `json:"execution_time"`
	MemoryUsed    float64 `json:"memory_used"`
	SubmittedAt   string  `json:"submitted_at"`
}
