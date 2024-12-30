package services

import (
	database "wow-bato-backend/internal"
	"wow-bato-backend/internal/models"
)

func CreateFeedback(newFeedback models.NewFeedback) error {
	db, err := database.ConnectDB()
    if err != nil {
        return err
    }

    feedback := models.Feedback{
        Content: newFeedback.Content,
        UserID: newFeedback.UserID,
        Role: newFeedback.Role,
    }

    result := db.Create(&feedback)

    return result.Error
}
