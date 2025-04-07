package handlers

import (
	"net/http"
	"strconv"
	"wow-bato-backend/internal/models"
	"wow-bato-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type FeedbackHandlers struct {
    svc *services.FeedbackService
}

func NewFeedbackHandlers(svc *services.FeedbackService) *FeedbackHandlers {
    return &FeedbackHandlers{svc: svc}
}

func (h *FeedbackHandlers) CreateFeedBack(c *gin.Context) {
    
    session := services.CheckAuthentication(c)

    var newFeedback models.NewFeedback
    services.BindJSON(c, &newFeedback)

    project_id := c.Param("projectID")
    user_id := session.Get("user_id").(uint)
    user_role := session.Get("user_role").(string)

    project_id_int, err := strconv.Atoi(project_id)
    services.CheckServiceError(c, err)

    feedback := models.CreateFeedback{
        Content: newFeedback.Content,
        Role: user_role,
        UserID: user_id,
        ProjectID: uint(project_id_int),
    }

    err = h.svc.CreateFeedback(feedback)
    services.CheckServiceError(c, err)

    c.IndentedJSON(http.StatusOK, gin.H{"message": "New feedback created"})
}

func (h *FeedbackHandlers) GetAllFeedbacks(c *gin.Context){
    
    services.CheckAuthentication(c)

    projectID := c.Param("projectID")

    feedbacks, err := h.svc.GetAllFeedback(projectID)
    services.CheckServiceError(c, err)

    c.IndentedJSON(http.StatusOK, gin.H{"feedbacks": feedbacks})

}

func (h *FeedbackHandlers) EditFeedback(c *gin.Context){
    
    services.CheckAuthentication(c)
    
    feedbackID := c.Param("feedbackID")

    var newFeedback models.NewFeedback
    services.BindJSON(c, &newFeedback)

    err := h.svc.EditFeedback(feedbackID, newFeedback)
    services.CheckServiceError(c, err)

    c.IndentedJSON(http.StatusOK, gin.H{"message": "Feedback edited"})
}

func (h *FeedbackHandlers) DeleteFeedback(c *gin.Context){

    services.CheckAuthentication(c)

    feedbackID := c.Param("feedbackID")

    err := h.svc.DeleteFeedback(feedbackID)
    services.CheckServiceError(c, err)

    c.IndentedJSON(http.StatusOK, gin.H{"message": "Feedback deleted"})
}