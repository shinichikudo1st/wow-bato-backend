package routes

import (
	"wow-bato-backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterFeedbackReplyRoutes registers feedback reply-related routes on the provided router group.
//
// This route contains:
// - POST /feedbackReply/create/:feedbackID: Handles creating a new feedback reply.
// - GET /feedbackReply/get/:feedbackID: Handles retrieving a list of feedback replies.
// - DELETE /feedbackReply/delete/:feedbackID: Handles deleting a feedback reply.
// - PUT /feedbackReply/edit/:replyID: Handles updating a feedback reply.
func RegisterFeedbackReplyRoutes(router *gin.RouterGroup) {
	feedbackReply := router.Group("/feedbackReply")
	{
		feedbackReply.POST("/create/:feedbackID", handlers.CreateFeedbackReply)
		feedbackReply.GET("/get/:feedbackID", handlers.GetAllReplies)
		feedbackReply.DELETE("/delete/:feedbackID", handlers.DeleteFeedbackReply)
		feedbackReply.PUT("/edit/:replyID", handlers.EditFeedbackReply)
	}
}
