package models

type NewBudgetCategory struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Barangay_ID uint   `json:"barangay_ID"`
}