package models

type NewFeedback struct {
    Content string `json:"content"`
    Role string `json:"role"`
    UserID uint `json:"user_id"`
    ProjectID uint `json:"project_id"`
}
