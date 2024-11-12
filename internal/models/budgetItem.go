package models

type NewBudgetItem struct {
	Name             string  `json:"name"`
	Amount_Allocated float64 `json:"amount_allocated"`
	Description      string  `json:"description"`
	Status           string  `json:"status"`
}

type UpdateStatus struct {
	Status string `json:"status"`
}