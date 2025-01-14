// Package services provides feedback-related business logic and operations for the application.
// It handles feedback management while maintaining separation of concerns from the database and presentation layers.
package services

import (
	"strconv"
	database "wow-bato-backend/internal"
	"wow-bato-backend/internal/models"
)

// CreateFeedback creates a new feedback
//
// This function performs the following operations:
// 1. Establishes a database connection
// 2. Creates a new feedback record in the database
// 3. Returns nil if successful, otherwise returns an error
//
// Parameters:
//   - newFeedback: models.CreateFeedback - The new feedback data
//
// Returns:
//   - error: Returns nil if successful, otherwise returns an error:
//   - Database connection errors
//   - Database creation errors
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

// GetAllFeedback retrieves all feedback for a specific project
//
// This function performs the following operations:
// 1. Establishes a database connection
// 2. Retrieves all feedback for the specified project ID
// 3. Returns nil if successful, otherwise returns an error
//
// Parameters:
//   - projectID: string - The ID of the project
//
// Returns:
//   - []models.GetAllFeedbacks: Returns the retrieved feedback
//   - error: Returns nil if successful, otherwise returns an error
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

// EditFeedback updates an existing feedback
//
// This function performs the following operations:
// 1. Establishes a database connection
// 2. Retrieves the feedback record based on the ID
// 3. Updates the feedback content
// 4. Returns nil if successful, otherwise returns an error
//
// Parameters:
//   - feedbackID: string - The ID of the feedback to be updated
//   - editedFeedback: models.NewFeedback - The updated feedback data
//
// Returns:
//   - error: Returns nil if successful, otherwise returns an error:
//   - Database connection errors
//   - Feedback not found errors
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

// DeleteFeedback deletes a feedback from the database
//
// This function performs the following operations:
// 1. Establishes a database connection
// 2. Retrieves the feedback record based on the ID
// 3. Deletes the feedback record from the database
// 4. Returns nil if successful, otherwise returns an error
//
// Parameters:
//   - feedbackID: string - The ID of the feedback to be deleted
//
// Returns:
//   - error: Returns nil if successful, otherwise returns an error:
//   - Database connection errors
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
