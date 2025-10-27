package controllers

import (
	"net/http"
	"project-z-backend/handlers"
	"project-z-backend/models"
	"project-z-backend/services"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var u models.User

	if err := c.ShouldBindJSON(&u); err != nil {
		handlers.HandleError(c, http.StatusBadRequest, err.Error())

	}

	user, err := services.Register(u)

	if err != nil {
		handlers.HandleError(c, http.StatusInternalServerError, err.Error())

	}

	handlers.HandleSuccess(c, http.StatusOK, user)
}

func Me(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		handlers.HandleError(c, http.StatusUnauthorized, "User ID not found in context")
		return
	}

	user, err := services.UserInfo(userID.(int64))
	if err != nil {
		handlers.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}

	handlers.HandleSuccess(c, http.StatusOK, user)
}

func Login(c *gin.Context) {
	var u models.User

	if err := c.ShouldBindJSON(&u); err != nil {
		handlers.HandleError(c, http.StatusBadRequest, err.Error())
	}

	token, err := services.Login(u)

	if err != nil {
		handlers.HandleError(c, http.StatusUnauthorized, err.Error())
	}

	handlers.HandleSuccess(c, http.StatusOK, gin.H{"token": token})
}
