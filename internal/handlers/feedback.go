// Package handlers provides HTTP request handlers for the wow-bato application.
// It implements handlers for feedback management operations, including:
//   - Feedback creation and initialization
//   - Feedback updates and modifications
//   - Feedback deletion and cleanup
//
// The package ensures proper authentication and authorization checks
// while maintaining data consistency across feedback operations.
package handlers

import (
	"net/http"
	"strconv"
	"wow-bato-backend/internal/models"
	"wow-bato-backend/internal/services"

	"github.com/gin-gonic/gin"
)

// CreateFeedBack handles the creation and initialization of a new feedback
//
// This handler performs the following operations:
//  1. Validates user authentication and authorization
//  2. Validates and binds the new feedback data
//  3. Delegates feedback creation to the services layer
//  4. Returns appropriate response based on operation result
//
// Security:
//  - Requires authenticated session
//  - Validates administrative privileges
//
// @Summary Create a new feedback
// @Description Creates a new feedback with the provided details
// @Tags Feedback
// @Accept json
// @Produce json
// @Param projectID path string true "Project ID"
// @Param feedback body models.NewFeedback true "Feedback details"
// @Success 200 {object} gin.H "Returns success message on feedback creation"
// @Failure 400 {object} gin.H "Returns error message on invalid request"
// @Failure 500 {object} gin.H "Returns error message on server error"
// @Router /feedback/{projectID} [post]
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

// GetAllFeedbacks handles the retrieval of all feedback for a specific project
//
// This handler performs the following operations:
//  1. Validates user authentication and authorization
//  2. Validates and binds the project ID
//  3. Delegates feedback retrieval to the services layer
//  4. Returns appropriate response based on operation result
//
// Security:
//  - Requires authenticated session
//  - Validates administrative privileges
//
// @Summary Get all feedback for a project
// @Description Retrieves all feedback for the specified project
// @Tags Feedback
// @Accept json
// @Produce json
// @Param projectID path string true "Project ID"
// @Success 200 {object} gin.H "Returns a list of feedback"
// @Failure 401 {object} gin.H "Returns error message on unauthorized access"
// @Failure 500 {object} gin.H "Returns error message on server error"
// @Router /feedback/{projectID} [get]
func GetAllFeedbacks(c *gin.Context){
    
    services.CheckAuthentication(c)

    projectID := c.Param("projectID")

    feedbacks, err := services.GetAllFeedback(projectID)
    services.CheckServiceError(c, err)

    c.IndentedJSON(http.StatusOK, gin.H{"feedbacks": feedbacks})

}

// EditFeedback handles the editing of an existing feedback
//
// This handler performs the following operations:
//  1. Validates user authentication and authorization
//  2. Validates and binds the feedback ID and new feedback data
//  3. Delegates feedback editing to the services layer
//  4. Returns appropriate response based on operation result
//
// Security:
//  - Requires authenticated session
//  - Validates administrative privileges
//
// @Summary Edit a feedback
// @Description Edits an existing feedback with the provided details
// @Tags Feedback
// @Accept json
// @Produce json
// @Param feedbackID path string true "Feedback ID"
// @Param feedback body models.NewFeedback true "Feedback details"
// @Success 200 {object} gin.H "Returns success message on feedback editing"
// @Failure 400 {object} gin.H "Returns error message on invalid request"
// @Failure 500 {object} gin.H "Returns error message on server error"
// @Router /feedback/{feedbackID} [put]
func EditFeedback(c *gin.Context){
    
    services.CheckAuthentication(c)
    
    feedbackID := c.Param("feedbackID")

    var newFeedback models.NewFeedback
    services.BindJSON(c, &newFeedback)

    err := services.EditFeedback(feedbackID, newFeedback)
    services.CheckServiceError(c, err)

    c.IndentedJSON(http.StatusOK, gin.H{"message": "Feedback edited"})
}

// DeleteFeedback handles the deletion of a feedback
//
// This handler performs the following operations:
//  1. Validates user authentication and authorization
//  2. Validates and binds the feedback ID
//  3. Delegates feedback deletion to the services layer
//  4. Returns appropriate response based on operation result
//
// Security:
//  - Requires authenticated session
//  - Validates administrative privileges
//
// @Summary Delete a feedback
// @Description Deletes an existing feedback with the provided ID
// @Tags Feedback
// @Accept json no body
// @Produce json
// @Param feedbackID path string true "Feedback ID"
// @Success 200 {object} gin.H "Returns success message on feedback deletion"
// @Failure 400 {object} gin.H "Returns error message on invalid request"
// @Failure 500 {object} gin.H "Returns error message on server error"
// @Router /feedback/{feedbackID} [delete]
func DeleteFeedback(c *gin.Context){

    services.CheckAuthentication(c)

    feedbackID := c.Param("feedbackID")

    err := services.DeleteFeedback(feedbackID)
    services.CheckServiceError(c, err)

    c.IndentedJSON(http.StatusOK, gin.H{"message": "Feedback deleted"})
}