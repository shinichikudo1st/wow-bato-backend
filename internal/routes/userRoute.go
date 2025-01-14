package routes

import (
	"wow-bato-backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterUserRoute registers user-related routes on the provided router group.
//
// This route contains:
// - POST /user/register: Handles user registration.
// - POST /user/login: Handles user login.
// - POST /user/logout: Handles user logout.
// - GET /user/checkAuth: Checks if the user is authenticated.
// - GET /user/profile: Retrieves the user's profile information.
func RegisterUserRoute(router *gin.RouterGroup) {
	user := router.Group("/user")
	{
		user.POST("/register", handlers.RegisterUser)
		user.POST("/login", handlers.LoginUser)
		user.POST("/logout", handlers.LogoutUser)
		user.GET("/checkAuth", handlers.CheckAuth)
		user.GET("/profile", handlers.GetUserProfile)
	}
}
