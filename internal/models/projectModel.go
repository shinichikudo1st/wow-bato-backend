package models

import "time"

type NewProject struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	StartDate   time.Time `json:"startDate"`
	Status string `json:"status"`
}

type NewProjectStatus struct {
    Status string `json:"status"`
    FlexDate time.Time `json:"flexdate"`
}
