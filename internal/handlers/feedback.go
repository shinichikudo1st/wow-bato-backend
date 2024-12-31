package handlers

import (
	"net/http"
	"wow-bato-backend/internal/models"
	"wow-bato-backend/internal/services"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

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

    err := services.CreateFeedback(newFeedback)
    if err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.IndentedJSON(http.StatusOK, gin.H{"message": "New feedback created"})
}

func GetAllFeedbacks(c *gin.Context){
    session := sessions.Default(c)

    if session.Get("authenticated") != true {
        c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Access Denied: Unauthorized"})
        return
    }

    projectID := c.Param("project_id")
    feedbacks, err := services.GetAllFeedback(projectID)
    if err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.IndentedJSON(http.StatusOK, gin.H{"feedbacks": feedbacks})

}
