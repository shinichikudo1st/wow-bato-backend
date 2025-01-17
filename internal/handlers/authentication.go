package handlers

import (
	"net/http"
	"wow-bato-backend/internal/models"
	"wow-bato-backend/internal/services"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Handler for user registration
// @Summary Register a new user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param registerUser body models.RegisterUser true "User registration details"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
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

// Handler for user login
// @Summary Login a user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param loginUser body models.LoginUser true "User login details"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
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

	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	session.Set("user_role", user.Role)
	session.Set("barangay_id", user.Barangay_ID)
	session.Set("barangay_name", user.Barangay_Name)
	session.Set("authenticated", true)
	

	if err := session.Save(); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	

	c.IndentedJSON(http.StatusOK, gin.H{"message": "User logged in successfully", "sessionStatus": session.Get("authenticated")})
}

// Handler for user logout
// @Summary Logout a user, clears the session
// @Tags Authentication
// @Accept json no body
// @Produce json
// @Success 200 {object} gin.H
// @Failure 500 {object} gin.H
func LogoutUser(c *gin.Context){

	session := sessions.Default(c)

	session.Clear()

	if err := session.Save(); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "User logged out successfully"})
}

// Handler for checking authentication
// @Summary Check authentication status
// @Tags Authentication
// @Accept json no body
// @Produce json
// @Success 200 {object} gin.H
// @Failure 500 {object} gin.H
func CheckAuth(c *gin.Context){
	session := sessions.Default(c)

	if session.Get("authenticated") != true {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"sessionStatus": session.Get("authenticated"), "role": session.Get("user_role"), 
	"user_id": session.Get("user_id"), "barangay_id": session.Get("barangay_id"), "barangay_name": session.Get("barangay_name")})
}

// Handler for getting user profile
// @Summary Get user profile, get user_id from session
// @Tags Authentication
// @Accept json no body
// @Produce json
// @Success 200 {object} gin.H
// @Failure 500 {object} gin.H
func GetUserProfile(c *gin.Context){
	session := sessions.Default(c)

	if session.Get("authenticated") != true {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Access Denied"})
		return
	}

	userID := session.Get("user_id")

	userProfile, err := services.GetUserProfile(userID.(uint))

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "User profile fetched successfully", "data": userProfile})
}

