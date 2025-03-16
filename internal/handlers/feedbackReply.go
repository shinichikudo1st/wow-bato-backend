// Package handlers provides HTTP request handlers for the wow-bato application.
// It implements handlers for feedback reply management operations, including:
//   - Feedback reply creation and initialization
//   - Feedback reply updates and modifications
//   - Feedback reply deletion and cleanup
//
// The package ensures proper authentication and authorization checks
// while maintaining data consistency across feedback reply operations.
package handlers

import (
	"net/http"
	"wow-bato-backend/internal/models"
	"wow-bato-backend/internal/services"

	"github.com/gin-gonic/gin"
)

// CreateFeedbackReply handles the creation and initialization of a new feedback reply
//
// This handler performs the following operations:
//  1. Validates user authentication and authorization
//  2. Validates and binds the feedback ID and new reply data
//  3. Delegates feedback reply creation to the services layer
//  4. Returns appropriate response based on operation result
//
// Security:
//  - Requires authenticated session
//  - Validates administrative privileges
//
// @Summary Create a new feedback reply
// @Description Creates a new feedback reply with the provided details
// @Tags Feedback Reply
// @Accept json
// @Produce json
// @Param feedbackID path string true "Feedback ID"
// @Param reply body models.Reply true "Reply details"
// @Success 200 {object} gin.H "Returns success message on feedback reply creation"
// @Failure 400 {object} gin.H "Returns error message on invalid request"
// @Failure 500 {object} gin.H "Returns error message on server error"
// @Router /feedback/{feedbackID}/reply [post]
func CreateFeedbackReply(c *gin.Context){
    
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

    err := services.CreateFeedbackReply(newReply)
    services.CheckServiceError(c, err)

    c.IndentedJSON(http.StatusOK, gin.H{"message": "Reply submitted"})
}

// GetAllReplies handles the retrieval of all replies for a specific feedback
// 
// This handler performs the following operations:
//  1. Validates user authentication and authorization
//  2. Validates and binds the feedback ID
//  3. Delegates feedback reply retrieval to the services layer
//  4. Returns appropriate response based on operation result
//
// Security:
//  - Requires authenticated session
//  - Validates administrative privileges
//
// @Summary Get all replies for a feedback
// @Description Retrieves all replies for a specific feedback
// @Tags Feedback Reply
// @Accept json
// @Produce json
// @Param feedbackID path string true "Feedback ID"
// @Success 200 {object} gin.H "Returns a list of feedback replies"
// @Failure 401 {object} gin.H "Returns error when user is not authenticated"
// @Failure 404 {object} gin.H "Returns error when feedback is not found"
// @Failure 500 {object} gin.H "Returns error when feedback reply retrieval fails"
// @Router /feedback/{feedbackID}/reply [get]
func GetAllReplies(c *gin.Context){
    
    services.CheckAuthentication(c)

    feedbackID := c.Param("feedbackID")

    replies, err := services.GetAllReplies(feedbackID)
    services.CheckServiceError(c, err)

    c.IndentedJSON(http.StatusOK, gin.H{"message": "Replies retrived", "data": replies})
}

// DeleteFeedbackReply handles the deletion of a feedback reply
// 
// This handler performs the following operations:
//  1. Validates user authentication and authorization
//  2. Validates and binds the reply ID
//  3. Delegates feedback reply deletion to the services layer
//  4. Returns appropriate response based on operation result
//
// Security:
//  - Requires authenticated session
//  - Validates administrative privileges
//
// @Summary Delete a feedback reply
// @Description Deletes an existing feedback reply with the provided ID
// @Tags Feedback Reply
// @Accept json
// @Produce json
// @Param replyID path string true "Reply ID"
// @Success 200 {object} gin.H "Returns success message on feedback reply deletion"
// @Failure 400 {object} gin.H "Returns error message on invalid request"
// @Failure 500 {object} gin.H "Returns error message on server error"
// @Router /feedback/reply/{replyID} [delete]
func DeleteFeedbackReply(c *gin.Context){
    
    services.CheckAuthentication(c)

    feedback_id := c.Param("feedbackID")

    err := services.DeleteFeedbackReply(feedback_id)
    services.CheckServiceError(c, err)

    c.IndentedJSON(http.StatusOK, gin.H{"message": "Reply deleted"})
}

// EditFeedbackReply handles the editing of a feedback reply
// 
// This handler performs the following operations:
//  1. Validates user authentication and authorization
//  2. Validates and binds the reply ID and new reply data
//  3. Delegates feedback reply editing to the services layer
//  4. Returns appropriate response based on operation result
//
// Security:
//  - Requires authenticated session
//  - Validates administrative privileges
//
// @Summary Edit a feedback reply
// @Description Edits an existing feedback reply with the provided details
// @Tags Feedback Reply
// @Accept json
// @Produce json
// @Param replyID path string true "Reply ID"
// @Param reply body models.Reply true "Reply details"
// @Success 200 {object} gin.H "Returns success message on feedback reply editing"
// @Failure 400 {object} gin.H "Returns error message on invalid request"
// @Failure 500 {object} gin.H "Returns error message on server error"
// @Router /feedback/reply/{replyID} [put]
func EditFeedbackReply(c *gin.Context){
    
    session := services.CheckAuthentication(c)

    var editReply models.EditReply
    services.BindJSON(c, &editReply)

    requestingID := editReply.UserID
    sessionID := session.Get("user_id").(uint)

    if requestingID != sessionID {
        c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Invalid User ID/ Unauthorized"})
    }

    replyID := c.Param("replyID")

    err := services.EditFeedbackReply(replyID, editReply.Content)
    services.CheckServiceError(c, err)

    c.IndentedJSON(http.StatusOK, gin.H{"message": "Reply Edited"})
}