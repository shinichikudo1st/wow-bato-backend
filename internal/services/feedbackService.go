package services

import (
	"strconv"
	database "wow-bato-backend/internal"
	"wow-bato-backend/internal/models"
)

func CreateFeedback(newFeedback models.CreateFeedback) error {
	db, err := database.ConnectDB()
    if err != nil {
        return err
    }

    feedback := models.Feedback{
        Content: newFeedback.Content,
        UserID: newFeedback.UserID,
        Role: newFeedback.Role,
        ProjectID: newFeedback.ProjectID,
    }

    result := db.Create(&feedback)

    return result.Error
}

func GetAllFeedback(projectID string)([]models.GetAllFeedbacks, error){
    db, err := database.ConnectDB()
    if err != nil {
        return []models.GetAllFeedbacks{}, err
    }

    projectid_int, err := strconv.Atoi(projectID)
    if err != nil {
        return []models.GetAllFeedbacks{}, err
    }
    
    var feedbacks []models.GetAllFeedbacks
    if err := db.Model(&models.Feedback{}).Where("project_id = ?", projectid_int).Select("id, content, role, project_id, user_id").Scan(&feedbacks).Error;
    err != nil {
        return []models.GetAllFeedbacks{}, err
    }

    return feedbacks, nil
}
