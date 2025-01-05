package services

import (
	"strconv"
	database "wow-bato-backend/internal"
	"wow-bato-backend/internal/models"
)

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
		Content: newReply.Content,
		FeedbackID: uint(feedbackID),
		UserID: newReply.UserID,
	}

	result := db.Create(&reply)

	return result.Error
}