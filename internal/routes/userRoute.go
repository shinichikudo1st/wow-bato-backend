package routes

import (
	"wow-bato-backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

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
