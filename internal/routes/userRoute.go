package routes

import (
	"wow-bato-backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoute(router *gin.RouterGroup){
	user := router.Group("/user")
	{
		user.POST("/registerUser", handlers.RegisterUser)
		user.POST("/loginUser", handlers.LoginUser)
		user.POST("/logoutUser", handlers.LogoutUser)
		user.GET("/checkAuth", handlers.CheckAuth)
	}
}
