package routes

import (
	"wow-bato-backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterFeedbackRoutes registers feedback-related routes on the provided router group.
//
// This route contains:
// - POST /feedback/create/:projectID: Handles creating a new feedback/comment.
// - GET /feedback/all/:projectID: Handles retrieving a list of feedbacks/comment.
// - PUT /feedback/update/:feedbackID: Handles updating a feedback/comment.
// - DELETE /feedback/delete/:feedbackID: Handles deleting a feedback/comment.
func RegisterFeedbackRoutes(router *gin.RouterGroup) {
	feedback := router.Group("/feedback")
	{
		feedback.POST("/create/:projectID", handlers.CreateFeedBack)
		feedback.GET("/all/:projectID", handlers.GetAllFeedbacks)
		feedback.PUT("/update/:feedbackID", handlers.EditFeedback)
		feedback.DELETE("/delete/:feedbackID", handlers.DeleteFeedback)
	}
}
