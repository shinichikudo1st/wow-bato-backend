// Package services provides feedback reply-related business logic and operations for the application.
// It handles feedback reply management while maintaining separation of concerns from the database and presentation layers.
package services

import (
	"strconv"
	database "wow-bato-backend/internal"
	"wow-bato-backend/internal/models"
)

// CreateFeedbackReply creates a new feedback reply in the database.
//
// This function performs the following operations:
// 1. Establishes a database connection
// 2. Creates a new feedback reply record in the database
// 3. Returns nil if successful, otherwise returns an error
//
// Parameters:
//   - newReply: models.NewFeedbackReply - The new feedback reply data
//
// Returns:
//   - error: Returns nil if successful, otherwise returns an error:
//   - Database connection errors
//   - Database creation errors
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

// GetAllReplies retrieves all replies for a specific feedback from the database.
//
// This function performs the following operations:
// 1. Establishes a database connection
// 2. Converts the feedbackID from string to int
// 3. Retrieves all replies for the specified feedback from the database
//
// Parameters:
//   - feedbackID: string - The ID of the feedback for which to retrieve replies
//
// Returns:
//   - []models.FeedbackReply: An array of feedback replies
//   - error: Returns nil if successful, otherwise returns an error
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
	if err := db.Where("id = ?", feedbackID_int).Find(&replies).Error; err != nil {
		return []models.FeedbackReply{}, err
	}

	return replies, nil
}

// DeleteFeedbackReply deletes a feedback reply from the database.
//
// This function performs the following operations:
// 1. Establishes a database connection
// 2. Converts the feedbackID from string to int
// 3. Deletes the feedback reply record from the database
//
// Parameters:
//   - feedbackID: string - The ID of the feedback reply to be deleted
//
// Returns:
//   - error: Returns nil if successful, otherwise returns an error
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

// EditFeedbackReply updates the content of a feedback reply in the database.
//
// This function performs the following operations:
// 1. Establishes a database connection
// 2. Converts the replyID from string to int
// 3. Retrieves the feedback reply record from the database that matches the replyID
// 4. Updates the content of the feedback reply record in the database
//
// Parameters:
//   - replyID: string - The ID of the feedback reply to be updated
//   - content: string - The new content of the feedback reply
//
// Returns:
//   - error: Returns nil if successful, otherwise returns an error
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
