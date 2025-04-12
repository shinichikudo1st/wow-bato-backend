package services

import (
	"strconv"
	"wow-bato-backend/internal/models"

	"gorm.io/gorm"
)

type FeedbackReplyService struct {
	db *gorm.DB
}

func NewFeedbackReplyService(db *gorm.DB) *FeedbackReplyService {
	return &FeedbackReplyService{db: db}
}

func (s *FeedbackReplyService) CreateFeedbackReply(newReply models.NewFeedbackReply) error {
	
	feedbackID, err := strconv.Atoi(newReply.FeedbackID)
	if err != nil {
		return err
	}

	reply := models.FeedbackReply{
		Content:    newReply.Content,
		FeedbackID: uint(feedbackID),
		UserID:     newReply.UserID,
	}

	result := s.db.Create(&reply)

	return result.Error
}

func (s *FeedbackReplyService) GetAllReplies(feedbackID string) ([]models.FeedbackReply, error) {

	feedbackID_int, err := strconv.Atoi(feedbackID)
	if err != nil {
		return []models.FeedbackReply{}, err
	}

	var replies []models.FeedbackReply
	if err := s.db.Where("feedback_ID = ?", feedbackID_int).Find(&replies).Error; err != nil {
		return []models.FeedbackReply{}, err
	}

	return replies, nil
}

func (s *FeedbackReplyService) DeleteFeedbackReply(feedbackID string) error {

	feedbackID_int, err := strconv.Atoi(feedbackID)
	if err != nil {
		return err
	}

	var reply models.FeedbackReply
	if err := s.db.Where("id = ?", feedbackID_int).Delete(&reply).Error; err != nil {
		return err
	}

	return nil
}

func (s *FeedbackReplyService) EditFeedbackReply(replyID string, content string) error {

	replyID_int, err := strconv.Atoi(replyID)
	if err != nil {
		return err
	}

	var reply models.FeedbackReply
	if err := s.db.Where("id = ?", replyID_int).First(&reply).Error; err != nil {
		return err
	}

	if content != "" {
		reply.Content = content
	}

	result := s.db.Save(&reply)

	return result.Error
}
