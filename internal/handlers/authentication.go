// Package handlers provides HTTP request handlers for the wow-bato application.
// It implements handlers for user authentication and session management, including:
//   - User registration and account creation
//   - User login with session management
//   - User logout and session cleanup
//   - Authentication status verification
//   - User profile management
//
// The package uses the Gin web framework for routing and request handling,
// and implements secure session management using gin-contrib/sessions.
package handlers

import (
	"net/http"
	"wow-bato-backend/internal/models"
	"wow-bato-backend/internal/services"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// RegisterUser handles the registration of new users in the system.
//
// This handler performs the following operations:
//  1. Validates and binds the incoming JSON request to RegisterUser model
//  2. Delegates user registration to the services layer
//  3. Returns appropriate HTTP status and response
//
// @Summary Register a new user in the system
// @Description Creates a new user account with the provided registration details
// @Tags Authentication
// @Accept json
// @Produce json
// @Param registerUser body models.RegisterUser true "User registration details including username, password, and role"
// @Success 200 {object} gin.H "Returns success message on successful registration"
// @Failure 400 {object} gin.H "Returns error when request validation fails"
// @Failure 500 {object} gin.H "Returns error when registration process fails"
// @Router /register [post]
func RegisterUser(c *gin.Context) {
	var registerUser models.RegisterUser

	services.BindJSON(c, &registerUser)

	err := services.RegisterUser(registerUser)

	services.CheckServiceError(c, err)

	c.IndentedJSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

// LoginUser handles user authentication and session creation.
//
// This handler performs the following operations:
//  1. Validates and binds the incoming JSON request to LoginUser model
//  2. Authenticates user credentials via the services layer
//  3. Creates and configures a new session for authenticated users
//  4. Stores essential user data in the session
//
// The session stores the following user information:
//  - user_id: Unique identifier of the authenticated user
//  - user_role: Role/permissions level of the user
//  - barangay_id: Associated barangay identifier
//  - barangay_name: Name of the associated barangay
//  - authenticated: Boolean flag indicating active session
//
// @Summary Authenticate user and create session
// @Description Validates user credentials and creates a new session
// @Tags Authentication
// @Accept json
// @Produce json
// @Param loginUser body models.LoginUser true "User credentials including username and password"
// @Success 200 {object} gin.H "Returns success message and session status"
// @Failure 400 {object} gin.H "Returns error when request validation fails"
// @Failure 500 {object} gin.H "Returns error when authentication or session creation fails"
// @Router /login [post]
func LoginUser(c *gin.Context){
	var loginUser models.LoginUser

	services.BindJSON(c, &loginUser)

	user, err := services.LoginUser(loginUser)
	services.CheckServiceError(c, err)

	session := sessions.Default(c)
	services.SetSession(session, user)

	err = session.Save()
	services.CheckServiceError(c, err)
	

    c.IndentedJSON(http.StatusOK, gin.H{"message": "User logged in successfully", "sessionStatus": session.Get("authenticated"), "role": session.Get("user_role")})

}

// LogoutUser handles user session termination and cleanup.
//
// This handler performs the following operations:
//  1. Retrieves the current session
//  2. Clears all session data
//  3. Persists the empty session to finalize logout
//
// @Summary Terminate user session
// @Description Clears the user session and logs out the user
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200 {object} gin.H "Returns success message on successful logout"
// @Failure 500 {object} gin.H "Returns error when session cleanup fails"
// @Router /logout [post]
func LogoutUser(c *gin.Context){

	session := sessions.Default(c)

	session.Clear()

	err := session.Save()
	services.CheckServiceError(c, err)

	c.IndentedJSON(http.StatusOK, gin.H{"message": "User logged out successfully"})
}

// CheckAuth verifies the authentication status of the current session.
//
// This handler performs the following operations:
//  1. Retrieves the current session
//  2. Checks the authentication status
//  3. Returns the current authentication state and user details
//
// The response includes:
//  - Authentication status
//  - User role if authenticated
//  - Associated barangay information if applicable
//
// @Summary Verify authentication status
// @Description Checks if the current session is authenticated and returns user details
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200 {object} gin.H "Returns authentication status and user details"
// @Failure 500 {object} gin.H "Returns error when session verification fails"
// @Router /check-auth [get]
func CheckAuth(c *gin.Context){
	

	session := services.CheckAuthentication(c)

	c.IndentedJSON(http.StatusOK, gin.H{"sessionStatus": session.Get("authenticated"), "role": session.Get("user_role"), 
	"user_id": session.Get("user_id"), "barangay_id": session.Get("barangay_id"), "barangay_name": session.Get("barangay_name")})
}

// GetUserProfile retrieves the profile information for the authenticated user.
//
// This handler performs the following operations:
//  1. Retrieves the current session
//  2. Extracts the user ID from the session
//  3. Fetches the user profile from the services layer
//
// The profile information includes:
//  - User basic details
//  - Associated barangay information
//  - Role and permissions
//
// @Summary Retrieve user profile
// @Description Gets the profile information for the currently authenticated user
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200 {object} gin.H "Returns user profile information"
// @Failure 500 {object} gin.H "Returns error when profile retrieval fails"
// @Router /profile [get]
func GetUserProfile(c *gin.Context){

	session := services.CheckAuthentication(c)

	userID := session.Get("user_id")

	userProfile, err := services.GetUserProfile(userID.(uint))

	services.CheckServiceError(c, err)

	c.IndentedJSON(http.StatusOK, gin.H{"message": "User profile fetched successfully", "data": userProfile})
}
