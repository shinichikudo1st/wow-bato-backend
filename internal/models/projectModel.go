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

type NewProjectStatus struct {
    Status string `json:"status"`
    FlexDate time.Time `json:"flexdate"`
}

type ProjectList struct {
    ID uint `json:"id"`
    Name string `json:"name"`
    StartDate string `json:"startDate"`
    EndDate string `json:"endDate"`
    Status string `json:"status"`
}