package handlers

import (
	"net/http"
	"wow-bato-backend/internal/models"
	"wow-bato-backend/internal/services"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func CreateFeedbackReply(c *gin.Context){
    session := sessions.Default(c)

    if session.Get("authenticated") != true {
        c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Access Denied: Unauthorized"})
        return
    }

    var reply models.Reply
    if err := c.ShouldBindJSON(&reply); err != nil {
        c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    feedback_id := c.Param("feedbackID")
    userID, ok := session.Get("user_id").(uint)
    if !ok {
        c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID"})
        return
    }

    newReply := models.NewFeedbackReply {
        Content: reply.Content,
        FeedbackID: feedback_id,
        UserID: userID,
    }

    err := services.CreateFeedbackReply(newReply)
    if err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.IndentedJSON(http.StatusOK, gin.H{"Message": "Reply submitted"})
}