package models

type Problem struct {
	ID                    int64      `json:"id"`
	Title                 string     `json:"title"`
	Description           string     `json:"description"`
	Difficulty            string     `json:"difficulty"`
	ExampleInput          string     `json:"example_input"`
	ExampleOutput         string     `json:"example_output"`
	CreatedAt             string     `json:"created_at"`
	TestCases             []TestCase `json:"test_cases"`
	TotalSubmissions      int64      `json:"total_submissions"`
	SuccessfulSubmissions int64      `json:"successful_submissions"`
	FailedSubmissions     int64      `json:"failed_submissions"`
	SuccessRate           float64    `json:"success_rate"`
	FailureRate           float64    `json:"failure_rate"`
}

type TestCase struct {
	ID             int64  `json:"id"`
	ProblemID      int64  `json:"problem_id"`
	Input          string `json:"input"`
	ExpectedOutput string `json:"expected_output"`
	IsSample       bool   `json:"is_sample"`
	CreatedAt      string `json:"created_at"`
}
