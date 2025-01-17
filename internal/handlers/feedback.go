package handlers

import (
	"net/http"
	"strconv"
	"wow-bato-backend/internal/models"
	"wow-bato-backend/internal/services"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Handler for creating feedback
// @Summary Create feedback
// @Tags Feedback
// @Accept json
// @Produce json
// @Param feedback body models.NewFeedback true "Feedback details"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
func CreateFeedBack(c *gin.Context) {
    session := sessions.Default(c)

    if session.Get("authenticated") != true {
        c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Access Denied: Unauthorized"})
        return
    }

    var newFeedback models.NewFeedback
    if err := c.ShouldBindJSON(&newFeedback); err != nil {
        c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    project_id := c.Param("projectID")
    user_id := session.Get("user_id").(uint)
    user_role := session.Get("user_role").(string)

    project_id_int, err := strconv.Atoi(project_id)
    if err != nil {
        c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    feedback := models.CreateFeedback{
        Content: newFeedback.Content,
        Role: user_role,
        UserID: user_id,
        ProjectID: uint(project_id_int),
    }

    err = services.CreateFeedback(feedback)
    if err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.IndentedJSON(http.StatusOK, gin.H{"message": "New feedback created"})
}

// Handler for getting all feedback
// @Summary Get all feedback
// @Tags Feedback
// @Accept json no body
// @Produce json
// @Param projectID path string true "Project ID"
// @Success 200 {object} gin.H
// @Failure 500 {object} gin.H
func GetAllFeedbacks(c *gin.Context){
    session := sessions.Default(c)

    if session.Get("authenticated") != true {
        c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Access Denied: Unauthorized"})
        return
    }

    projectID := c.Param("projectID")
    feedbacks, err := services.GetAllFeedback(projectID)
    if err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.IndentedJSON(http.StatusOK, gin.H{"feedbacks": feedbacks})

}

// Handler for editing feedback
// @Summary Edit feedback
// @Tags Feedback
// @Accept json
// @Produce json
// @Param feedback body models.NewFeedback true "Feedback details"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
func EditFeedback(c *gin.Context){
    session := sessions.Default(c)

    if session.Get("authenticated") != true {
        c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Access Denied: Unauthorized"})
        return
    }

    var newFeedback models.NewFeedback
    if err := c.ShouldBindJSON(&newFeedback); err != nil {
        c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    feedbackID := c.Param("feedbackID")
    err := services.EditFeedback(feedbackID, newFeedback)
    if err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.IndentedJSON(http.StatusOK, gin.H{"message": "Feedback edited"})
}

// Handler for deleting feedback
// @Summary Delete feedback
// @Tags Feedback
// @Accept json no body
// @Produce json
// @Param feedbackID path string true "Feedback ID"
// @Success 200 {object} gin.H
// @Failure 500 {object} gin.H
func DeleteFeedback(c *gin.Context){
    session := sessions.Default(c)

    if session.Get("authenticated") != true {
        c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Access Denied: Unauthorized"})
        return
    }

    feedbackID := c.Param("feedbackID")
    err := services.DeleteFeedback(feedbackID)
    if err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.IndentedJSON(http.StatusOK, gin.H{"message": "Feedback deleted"})
}