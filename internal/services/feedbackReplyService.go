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

func GetAllReplies(feedbackID string)([]models.FeedbackReply, error){
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