package controllers

import (
	"log"
	"net/http"
	"net/mail"
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

	if _, err := mail.ParseAddress(u.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email address"})
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
