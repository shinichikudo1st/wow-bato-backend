package routes

import (
	"wow-bato-backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterFeedbackRoutes(router *gin.RouterGroup, handlers *handlers.FeedbackHandlers) {
	feedback := router.Group("/feedback")
	{
		feedback.POST("/create/:projectID", handlers.CreateFeedBack)
		feedback.GET("/all/:projectID", handlers.GetAllFeedbacks)
		feedback.PUT("/update/:feedbackID", handlers.EditFeedback)
		feedback.DELETE("/delete/:feedbackID", handlers.DeleteFeedback)
	}
}
