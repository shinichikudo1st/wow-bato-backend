package models

type NewBudgetItem struct {
	Name             string  `json:"name"`
	Amount_Allocated float64 `json:"amount_allocated"`
	Amount_Spent     float64 `json:"amount_spent"`
	Description      string  `json:"description"`
	Status           string  `json:"status"`
}