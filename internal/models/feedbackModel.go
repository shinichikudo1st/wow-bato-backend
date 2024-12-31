package models

type NewFeedback struct {
    Content string `json:"content"`
    Role string `json:"role"`
    UserID uint `json:"user_id"`
    ProjectID uint `json:"project_id"`
}

type GetAllFeedbacks struct {
    ID uint `json:"feedback_id"`
    Content string `json:"content"`
    Role string `json:"role"`
    UserID uint `json:"user_id"`
    ProjectID uint `json:"project_id"`
}
