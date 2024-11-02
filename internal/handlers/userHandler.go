package handlers

import (
	"net/http"
	"wow-bato-backend/internal/models"
	"wow-bato-backend/internal/services"

	"github.com/gin-gonic/gin"
)

func RegisterUser(c *gin.Context) {
	var registerUser models.RegisterUser

	if err := c.ShouldBindJSON(&registerUser); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := services.RegisterUser(registerUser)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func LoginUser(c *gin.Context){
	var loginUser models.LoginUser

	if err := c.ShouldBindJSON(&loginUser); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := services.LoginUser(loginUser)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "User logged in successfully", "user": user})
}