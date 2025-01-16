package models

import "time"

type NewProject struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	StartDate   string `json:"startDate"`
    EndDate     string `json:"endDate"`
	Status string `json:"status"`
}

type UpdateProject struct {
    Name string `json:"name"`
    Description string `json:"description"`
}

// For displaying Project Status at the Client Side
type NewProjectStatus struct {
    Status string `json:"status"`
    FlexDate time.Time `json:"flexdate"`
}

// Projects are displayed in projectList.jsx
type ProjectList struct {
    ID uint `json:"id"`
    Name string `json:"name"`
    StartDate string `json:"startDate"`
    EndDate string `json:"endDate"`
    Status string `json:"status"`
}
