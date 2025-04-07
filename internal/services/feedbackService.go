package services

import (
	"strconv"
	"wow-bato-backend/internal/models"

	"gorm.io/gorm"
)

type FeedbackService struct {
	db *gorm.DB
}

func NewFeedbackService(db *gorm.DB) *FeedbackService {
	return &FeedbackService{db: db}
}

func (s *FeedbackService) CreateFeedback(newFeedback models.CreateFeedback) error {

	feedback := models.Feedback{
		Content:   newFeedback.Content,
		UserID:    newFeedback.UserID,
		Role:      newFeedback.Role,
		ProjectID: newFeedback.ProjectID,
	}

	result := s.db.Create(&feedback)

	return result.Error
}

func (s *FeedbackService) GetAllFeedback(projectID string) ([]models.GetAllFeedbacks, error) {

	projectid_int, err := strconv.Atoi(projectID)
	if err != nil {
		return []models.GetAllFeedbacks{}, err
	}

	var feedbacks []models.GetAllFeedbacks
	if err := s.db.Model(&models.Feedback{}).Where("project_id = ?", projectid_int).Select("id, content, role, project_id, user_id").Scan(&feedbacks).Error; err != nil {
		return []models.GetAllFeedbacks{}, err
	}

	var user_id_list []uint
	for _, feedback := range feedbacks {
		user_id_list = append(user_id_list, feedback.UserID)
	}

	var users []models.FeedbackUser
	if err := s.db.Model(&models.User{}).Where("id IN (?)", user_id_list).Select("id, first_name, last_name").Scan(&users).Error; err != nil {
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

func (s *FeedbackService) EditFeedback(feedbackID string, editedFeedback models.NewFeedback) error {

	feedbackID_int, err := strconv.Atoi(feedbackID)
	if err != nil {
		return err
	}

	var feedback models.Feedback
	if err := s.db.Where("id = ?", feedbackID_int).First(&feedback).Error; err != nil {
		return err
	}

	feedback.Content = editedFeedback.Content

	result := s.db.Save(&feedback)

	return result.Error
}

func (s *FeedbackService) DeleteFeedback(feedbackID string) error {

	feedbackID_int, err := strconv.Atoi(feedbackID)
	if err != nil {
		return err
	}

	var feedback models.Feedback
	if err := s.db.Where("id = ?", feedbackID_int).Delete(&feedback).Error; err != nil {
		return err
	}

	return nil
}
