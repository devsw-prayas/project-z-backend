package controllers

import (
	"log"
	"net/http"
	"project-z-backend/models"
	"project-z-backend/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetProblems(c *gin.Context) {
	log.Println("GetProblems controller called")

	problems, err := services.GetProblems()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": problems,
	})
}

func CreateProblem(c *gin.Context) {
	log.Println("CreateProblem controller called")

	var p models.Problem
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	problem, err := services.CreateProblem(p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Problem created successfully",
		"data":    problem,
	})
}

func GetProblemByID(c *gin.Context) {
	log.Println("GetProblemByID controller called")

	problemIDStr := c.Param("problem_id")
	problemID, err := strconv.ParseInt(problemIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid problem ID"})
		return
	}

	problem, err := services.GetProblemByID(problemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if problem == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Problem not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": problem,
	})
}

func SubmitSolution(c *gin.Context) {
	log.Println("SubmitSolution controller called")

	var submission models.Submission
	if err := c.ShouldBindJSON(&submission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID := c.GetInt64("user_id")
	if userID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id not found in context"})
		return
	}

	submission.UserID = userID
	submission.Status = "Queued"

	result, err := services.CreateSubmission(submission)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Submission created successfully",
		"data":    result,
	})
}

func GetSubmission(c *gin.Context) {
	log.Println("GetSubmission controller called")

	submissionIDStr := c.Param("submission_id")
	submissionID, err := strconv.ParseInt(submissionIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid submission ID"})
		return
	}

	submission, err := services.GetSubmissionByID(submissionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if submission == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Submission not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": submission,
	})
}
