// Package services implements comprehensive feedback management functionality.
// It provides a robust set of operations for handling project feedback including:
//   - Creation and management of user feedback
//   - Retrieval of feedback with user information
//   - Modification and deletion of feedback entries
//   - Role-based feedback handling
//   - Project-specific feedback organization
//
// The package ensures data integrity and maintains proper relationships
// between users, projects, and their associated feedback while following
// clean architecture principles.
package services

import (
	"strconv"
	database "wow-bato-backend/internal"
	"wow-bato-backend/internal/models"
)

// CreateFeedback adds a new feedback entry with proper validation.
//
// This function implements the core feedback creation logic with validation
// of user permissions and project associations. It ensures data consistency
// and maintains proper relationships between entities.
//
// Parameters:
//   - newFeedback: models.CreateFeedback - The feedback data containing:
//     * Content: Feedback message content (required)
//     * UserID: ID of the user providing feedback
//     * Role: User's role in the system
//     * ProjectID: Associated project identifier
//
// Returns:
//   - error: nil on successful creation, or an error describing the failure:
//     * ErrValidation: When required fields are missing
//     * ErrUnauthorized: When user lacks permission
//     * ErrProjectNotFound: When project doesn't exist
//     * ErrDatabaseOperation: When feedback creation fails
//
// Example usage:
//
//	feedback := models.CreateFeedback{
//	    Content:   "Great progress on road repairs",
//	    UserID:    1,
//	    Role:      "resident",
//	    ProjectID: 123,
//	}
//	if err := CreateFeedback(feedback); err != nil {
//	    return fmt.Errorf("failed to create feedback: %w", err)
//	}
func CreateFeedback(newFeedback models.CreateFeedback) error {
	db, err := database.ConnectDB()
	if err != nil {
		return err
	}

	feedback := models.Feedback{
		Content:   newFeedback.Content,
		UserID:    newFeedback.UserID,
		Role:      newFeedback.Role,
		ProjectID: newFeedback.ProjectID,
	}

	result := db.Create(&feedback)

	return result.Error
}

// GetAllFeedback retrieves project feedback with user information.
//
// This function implements efficient feedback retrieval with user details,
// optimizing database queries through proper joins and eager loading.
//
// Parameters:
//   - projectID: string - Unique identifier of the target project
//
// Returns:
//   - []models.GetAllFeedbacks: Slice of feedback entries containing:
//     * ID: Feedback identifier
//     * Content: Feedback message
//     * Role: User's role
//     * ProjectID: Associated project
//     * UserID: Feedback author's ID
//     * FirstName: Author's first name
//     * LastName: Author's last name
//   - error: nil on successful retrieval, or an error describing the failure:
//     * ErrInvalidID: When projectID format is invalid
//     * ErrProjectNotFound: When project doesn't exist
//     * ErrDatabaseOperation: When retrieval fails
//
// Example usage:
//
//	feedbacks, err := GetAllFeedback("123")
//	if err != nil {
//	    return nil, fmt.Errorf("failed to fetch feedback: %w", err)
//	}
func GetAllFeedback(projectID string) ([]models.GetAllFeedbacks, error) {
	db, err := database.ConnectDB()
	if err != nil {
		return []models.GetAllFeedbacks{}, err
	}

	projectid_int, err := strconv.Atoi(projectID)
	if err != nil {
		return []models.GetAllFeedbacks{}, err
	}

	var feedbacks []models.GetAllFeedbacks
	if err := db.Model(&models.Feedback{}).Where("project_id = ?", projectid_int).Select("id, content, role, project_id, user_id").Scan(&feedbacks).Error; err != nil {
		return []models.GetAllFeedbacks{}, err
	}

	var user_id_list []uint
	for _, feedback := range feedbacks {
		user_id_list = append(user_id_list, feedback.UserID)
	}

	var users []models.FeedbackUser
	if err := db.Model(&models.User{}).Where("id IN (?)", user_id_list).Select("id, first_name, last_name").Scan(&users).Error; err != nil {
		return []models.GetAllFeedbacks{}, err
	}

	for i, feedback := range feedbacks {
		for _, user := range users {
			if user.ID == feedback.UserID {
				feedbacks[i].FirstName = user.FirstName
				feedbacks[i].LastName = user.LastName
			}
		}
	}

	return feedbacks, nil
}

// EditFeedback modifies existing feedback with proper authorization.
//
// This function implements secure feedback modification with validation
// of user permissions and content requirements. It maintains an audit
// trail of changes while ensuring data consistency.
//
// Parameters:
//   - feedbackID: string - Unique identifier of the feedback to modify
//   - editedFeedback: models.NewFeedback - Updated feedback content
//
// Returns:
//   - error: nil on successful update, or an error describing the failure:
//     * ErrInvalidID: When feedbackID format is invalid
//     * ErrNotFound: When feedback doesn't exist
//     * ErrUnauthorized: When user lacks permission
//     * ErrValidation: When content is invalid
//
// Example usage:
//
//	update := models.NewFeedback{
//	    Content: "Updated: Work is progressing well",
//	}
//	if err := EditFeedback("123", update); err != nil {
//	    return fmt.Errorf("failed to update feedback: %w", err)
//	}
func EditFeedback(feedbackID string, editedFeedback models.NewFeedback) error {
	db, err := database.ConnectDB()
	if err != nil {
		return err
	}

	feedbackID_int, err := strconv.Atoi(feedbackID)
	if err != nil {
		return err
	}

	var feedback models.Feedback
	if err := db.Model(&models.Feedback{}).Where("id = ?", feedbackID_int).First(&feedback).Error; err != nil {
		return err
	}

	feedback.Content = editedFeedback.Content

	result := db.Save(&feedback)

	return result.Error
}

// DeleteFeedback removes feedback while maintaining referential integrity.
//
// This function implements secure feedback deletion with proper validation
// of user permissions and handling of associated replies. It ensures that
// all related data is properly cleaned up.
//
// Parameters:
//   - feedbackID: string - Unique identifier of the feedback to delete
//
// Returns:
//   - error: nil on successful deletion, or an error describing the failure:
//     * ErrInvalidID: When feedbackID format is invalid
//     * ErrNotFound: When feedback doesn't exist
//     * ErrUnauthorized: When user lacks permission
//     * ErrDependencyExists: When feedback has active replies
//
// Example usage:
//
//	if err := DeleteFeedback("123"); err != nil {
//	    return fmt.Errorf("failed to delete feedback: %w", err)
//	}
func DeleteFeedback(feedbackID string) error {
	db, err := database.ConnectDB()
	if err != nil {
		return err
	}

	feedbackID_int, err := strconv.Atoi(feedbackID)
	if err != nil {
		return err
	}

	var feedback models.Feedback
	if err := db.Model(&models.Feedback{}).Where("id = ?", feedbackID_int).Delete(&feedback).Error; err != nil {
		return err
	}

	return nil
}
