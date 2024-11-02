package routes

import (
	"wow-bato-backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoute(router *gin.Engine){
	router.POST("/registerUser", handlers.RegisterUser)
}

func LoginUserRoute(router *gin.Engine){
	router.POST("/loginUser", handlers.LoginUser)
}