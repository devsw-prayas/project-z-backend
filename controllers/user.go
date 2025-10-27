package controllers

import (
	"log"
	"net/http"
	"project-z-backend/models"
	"project-z-backend/services"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	log.Println("Register controller called")

	var u models.User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := services.Register(u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func UserInfo(c *gin.Context) {
	log.Println("UserInfo controller called")
	var u models.User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := services.UserInfo(u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func Login(c *gin.Context) {
	log.Println("Login controller called")
	var u models.User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := services.Login(u)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func GetUserStats(c *gin.Context) {
	log.Println("GetUserStats controller called")

	userID := c.GetInt64("user_id")
	if userID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id not found in context"})
		return
	}

	user, err := services.GetUserWithStats(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}

func GetUserSubmissions(c *gin.Context) {
	log.Println("GetUserSubmissions controller called")

	userID := c.GetInt64("user_id")
	if userID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id not found in context"})
		return
	}

	submissions, err := services.GetUserSubmissions(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": submissions,
	})
}
