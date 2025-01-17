package handlers

import (
	"net/http"
	"wow-bato-backend/internal/models"
	"wow-bato-backend/internal/services"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Handler for creating feedback reply
// @Summary Create feedback reply
// @Tags Feedback Reply
// @Accept json
// @Produce json
// @Param reply body models.Reply true "Reply details"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
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

    c.IndentedJSON(http.StatusOK, gin.H{"message": "Reply submitted"})
}

// Handler for getting all replies
// @Summary Get all replies
// @Tags Feedback Reply
// @Accept json no body
// @Produce json
// @Param feedbackID path string true "Feedback ID"
// @Success 200 {object} gin.H
// @Failure 500 {object} gin.H
func GetAllReplies(c *gin.Context){
    session := sessions.Default(c)

    if session.Get("authenticated") != true {
        c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Access Denied: Unauthorized"})
        return
    }

    feedbackID := c.Param("feedbackID")

    replies, err := services.GetAllReplies(feedbackID)
    if err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.IndentedJSON(http.StatusOK, gin.H{"message": "Replies retrived", "data": replies})
}

// Handler for deleting feedback reply
// @Summary Delete feedback reply
// @Tags Feedback Reply
// @Accept json no body
// @Produce json
// @Param feedbackID path string true "Feedback ID"
// @Success 200 {object} gin.H
// @Failure 500 {object} gin.H
func DeleteFeedbackReply(c *gin.Context){
    session := sessions.Default(c)

    if session.Get("authenticated") != true {
        c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Access Denied: Unauthorized"})
        return
    }

    feedback_id := c.Param("feedbackID")

    err := services.DeleteFeedbackReply(feedback_id)
    if err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }


    c.IndentedJSON(http.StatusOK, gin.H{"message": "Reply deleted"})
}

// Handler for editing feedback reply
// @Summary Edit feedback reply
// @Tags Feedback Reply
// @Accept json
// @Produce json
// @Param reply body models.Reply true "Reply details"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
func EditFeedbackReply(c *gin.Context){
    session := sessions.Default(c)

    if session.Get("authenticated") != true {
        c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Access Denied: Unauthorized"})
        return
    }

    var editReply models.Reply
    if err := c.ShouldBindJSON(&editReply); err != nil {
        c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    replyID := c.Param("replyID")

    err := services.EditFeedbackReply(replyID, editReply.Content)
    if err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.IndentedJSON(http.StatusOK, gin.H{"message": "Reply Edited"})
}