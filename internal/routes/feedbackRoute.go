package routes

import (
	"wow-bato-backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterFeedbackRoutes(router *gin.RouterGroup) {
	feedback := router.Group("/feedback")
	{
        feedback.POST("/create/:projectID", handlers.CreateFeedBack)
        feedback.GET("/all/:projectID", handlers.GetAllFeedbacks)
	}
}
