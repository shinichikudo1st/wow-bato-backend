package routes

import "github.com/gin-gonic/gin"

func RegisterFeedbackReplyRoutes(router *gin.RouterGroup) {
	feedbackReply := router.Group("/feedbackReply")
	{
		feedbackReply.POST("/create/:feedbackID")
		feedbackReply.GET("/get/:feedbackID")
	}
}