package routes

import (
	"wow-bato-backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoute(router *gin.Engine){
	router.POST("/registerUser", handlers.RegisterUser)
}