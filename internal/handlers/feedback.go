package handlers

import (
	"net/http"
	"strconv"
	"wow-bato-backend/internal/models"
	"wow-bato-backend/internal/services"

	"github.com/gin-gonic/gin"
)

func CreateFeedBack(c *gin.Context) {
    
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

    err = services.CreateFeedback(feedback)
    services.CheckServiceError(c, err)

    c.IndentedJSON(http.StatusOK, gin.H{"message": "New feedback created"})
}

func GetAllFeedbacks(c *gin.Context){
    
    services.CheckAuthentication(c)

    projectID := c.Param("projectID")

    feedbacks, err := services.GetAllFeedback(projectID)
    services.CheckServiceError(c, err)

    c.IndentedJSON(http.StatusOK, gin.H{"feedbacks": feedbacks})

}

func EditFeedback(c *gin.Context){
    
    services.CheckAuthentication(c)
    
    feedbackID := c.Param("feedbackID")

    var newFeedback models.NewFeedback
    services.BindJSON(c, &newFeedback)

    err := services.EditFeedback(feedbackID, newFeedback)
    services.CheckServiceError(c, err)

    c.IndentedJSON(http.StatusOK, gin.H{"message": "Feedback edited"})
}

func DeleteFeedback(c *gin.Context){

    services.CheckAuthentication(c)

    feedbackID := c.Param("feedbackID")

    err := services.DeleteFeedback(feedbackID)
    services.CheckServiceError(c, err)

    c.IndentedJSON(http.StatusOK, gin.H{"message": "Feedback deleted"})
}