package handlers

import (
	"net/http"
	"strconv"
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
