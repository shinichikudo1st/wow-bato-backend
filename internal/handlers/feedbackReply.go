package handlers

import (
	"net/http"
	"wow-bato-backend/internal/models"
	"wow-bato-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type FeedbackReplyHandlers struct {
    svc *services.FeedbackReplyService
}

func NewFeedbackReplyHandlers(svc *services.FeedbackReplyService) *FeedbackReplyHandlers {
    return &FeedbackReplyHandlers{svc: svc}
}

func (h *FeedbackReplyHandlers) CreateFeedbackReply(c *gin.Context){
    
    session := services.CheckAuthentication(c)

    var reply models.Reply
    services.BindJSON(c, &reply)

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

    err := h.svc.CreateFeedbackReply(newReply)
    services.CheckServiceError(c, err)

    c.IndentedJSON(http.StatusOK, gin.H{"message": "Reply submitted"})
}

func (h *FeedbackReplyHandlers) GetAllReplies(c *gin.Context){
    
    services.CheckAuthentication(c)

    feedbackID := c.Param("feedbackID")

    replies, err := h.svc.GetAllReplies(feedbackID)
    services.CheckServiceError(c, err)

    c.IndentedJSON(http.StatusOK, gin.H{"message": "Replies retrived", "data": replies})
}

func (h *FeedbackReplyHandlers) DeleteFeedbackReply(c *gin.Context){
    
    services.CheckAuthentication(c)

    feedback_id := c.Param("feedbackID")

    err := h.svc.DeleteFeedbackReply(feedback_id)
    services.CheckServiceError(c, err)

    c.IndentedJSON(http.StatusOK, gin.H{"message": "Reply deleted"})
}

func (h *FeedbackReplyHandlers) EditFeedbackReply(c *gin.Context){
    
    session := services.CheckAuthentication(c)

    var editReply models.EditReply
    services.BindJSON(c, &editReply)

    requestingID := editReply.UserID
    sessionID := session.Get("user_id").(uint)

    if requestingID != sessionID {
        c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Invalid User ID/ Unauthorized"})
    }

    replyID := c.Param("replyID")

    err := h.svc.EditFeedbackReply(replyID, editReply.Content)
    services.CheckServiceError(c, err)

    c.IndentedJSON(http.StatusOK, gin.H{"message": "Reply Edited"})
}