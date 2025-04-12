package routes

import (
	"wow-bato-backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterFeedbackReplyRoutes(router *gin.RouterGroup, handlers *handlers.FeedbackReplyHandlers ) {
	feedbackReply := router.Group("/feedbackReply")
	{
		feedbackReply.POST("/create/:feedbackID", handlers.CreateFeedbackReply)
		feedbackReply.GET("/get/:feedbackID", handlers.GetAllReplies)
		feedbackReply.DELETE("/delete/:feedbackID", handlers.DeleteFeedbackReply)
		feedbackReply.PUT("/edit/:replyID", handlers.EditFeedbackReply)
	}
}
