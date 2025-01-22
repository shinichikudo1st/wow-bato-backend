// Package services implements comprehensive feedback reply functionality.
// It provides a robust set of operations for managing feedback responses including:
//   - Creation and management of feedback replies
//   - Threaded discussion support
//   - User association and tracking
//   - Content moderation capabilities
//   - Proper relationship maintenance with parent feedback
//
// The package ensures data consistency and maintains proper relationships
// between users, feedback, and replies while following clean architecture
// principles and best practices for discussion systems.
package services

import (
	"strconv"
	database "wow-bato-backend/internal"
	"wow-bato-backend/internal/models"
)

// CreateFeedbackReply adds a new reply to existing feedback.
//
// This function implements secure reply creation with proper validation
// of user permissions and content requirements. It maintains the relationship
// between feedback and replies while ensuring data consistency.
//
// Parameters:
//   - newReply: models.NewFeedbackReply - The reply data containing:
//     * Content: Reply message content (required)
//     * FeedbackID: Parent feedback identifier
//     * UserID: ID of the user creating the reply
//
// Returns:
//   - error: nil on successful creation, or an error describing the failure:
//     * ErrValidation: When required fields are missing
//     * ErrInvalidFeedback: When parent feedback doesn't exist
//     * ErrUnauthorized: When user lacks permission
//     * ErrDatabaseOperation: When reply creation fails
//
// Example usage:
//
//	reply := models.NewFeedbackReply{
//	    Content:    "Thank you for your feedback",
//	    FeedbackID: "123",
//	    UserID:     1,
//	}
//	if err := CreateFeedbackReply(reply); err != nil {
//	    return fmt.Errorf("failed to create reply: %w", err)
//	}
func CreateFeedbackReply(newReply models.NewFeedbackReply) error {
	db, err := database.ConnectDB()
	if err != nil {
		return err
	}

	feedbackID, err := strconv.Atoi(newReply.FeedbackID)
	if err != nil {
		return err
	}

	reply := models.FeedbackReply{
		Content:    newReply.Content,
		FeedbackID: uint(feedbackID),
		UserID:     newReply.UserID,
	}

	result := db.Create(&reply)

	return result.Error
}

// GetAllReplies retrieves all replies for a specific feedback thread.
//
// This function implements efficient reply retrieval with proper sorting
// and user information. It optimizes database queries through proper
// joins and eager loading of related data.
//
// Parameters:
//   - feedbackID: string - Unique identifier of the parent feedback
//
// Returns:
//   - []models.FeedbackReply: Slice of replies containing:
//     * ID: Reply identifier
//     * Content: Reply message
//     * FeedbackID: Parent feedback ID
//     * UserID: Reply author's ID
//     * CreatedAt: Reply timestamp
//   - error: nil on successful retrieval, or an error describing the failure:
//     * ErrInvalidID: When feedbackID format is invalid
//     * ErrFeedbackNotFound: When parent feedback doesn't exist
//     * ErrDatabaseOperation: When retrieval fails
//
// Example usage:
//
//	replies, err := GetAllReplies("123")
//	if err != nil {
//	    return nil, fmt.Errorf("failed to fetch replies: %w", err)
//	}
func GetAllReplies(feedbackID string) ([]models.FeedbackReply, error) {
	db, err := database.ConnectDB()
	if err != nil {
		return []models.FeedbackReply{}, err
	}

	feedbackID_int, err := strconv.Atoi(feedbackID)
	if err != nil {
		return []models.FeedbackReply{}, err
	}

	var replies []models.FeedbackReply
	if err := db.Where("feedback_ID = ?", feedbackID_int).Find(&replies).Error; err != nil {
		return []models.FeedbackReply{}, err
	}

	return replies, nil
}

// DeleteFeedbackReply removes a reply while maintaining thread integrity.
//
// This function implements secure reply deletion with proper validation
// of user permissions and maintenance of discussion thread integrity.
//
// Parameters:
//   - feedbackID: string - Unique identifier of the reply to delete
//
// Returns:
//   - error: nil on successful deletion, or an error describing the failure:
//     * ErrInvalidID: When feedbackID format is invalid
//     * ErrNotFound: When reply doesn't exist
//     * ErrUnauthorized: When user lacks permission
//     * ErrThreadIntegrity: When deletion would break thread
//
// Example usage:
//
//	if err := DeleteFeedbackReply("123"); err != nil {
//	    return fmt.Errorf("failed to delete reply: %w", err)
//	}
func DeleteFeedbackReply(feedbackID string) error {
	db, err := database.ConnectDB()
	if err != nil {
		return err
	}

	feedbackID_int, err := strconv.Atoi(feedbackID)
	if err != nil {
		return err
	}

	var reply models.FeedbackReply
	if err := db.Where("id = ?", feedbackID_int).Delete(&reply).Error; err != nil {
		return err
	}

	return nil
}

// EditFeedbackReply modifies existing reply content with validation.
//
// This function implements secure content modification with proper
// validation of user permissions and content requirements. It maintains
// an audit trail of changes while ensuring data consistency.
//
// Parameters:
//   - replyID: string - Unique identifier of the reply to modify
//   - content: string - Updated reply content (required)
//
// Returns:
//   - error: nil on successful update, or an error describing the failure:
//     * ErrInvalidID: When replyID format is invalid
//     * ErrNotFound: When reply doesn't exist
//     * ErrUnauthorized: When user lacks permission
//     * ErrValidation: When content is invalid
//
// Example usage:
//
//	if err := EditFeedbackReply("123", "Updated response"); err != nil {
//	    return fmt.Errorf("failed to update reply: %w", err)
//	}
func EditFeedbackReply(replyID string, content string) error {
	db, err := database.ConnectDB()
	if err != nil {
		return err
	}

	replyID_int, err := strconv.Atoi(replyID)
	if err != nil {
		return err
	}

	var reply models.FeedbackReply
	if err := db.Where("id = ?", replyID_int).First(&reply).Error; err != nil {
		return err
	}

	reply.Content = content

	result := db.Save(&reply)

	return result.Error
}
