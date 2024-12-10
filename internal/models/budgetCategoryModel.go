package models

type NewBudgetCategory struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Barangay_ID uint   `json:"barangay_ID"`
}

type UpdateBudgetCategory struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type BudgetCategoryResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Barangay_ID uint   `json:"barangay_ID"`
}

type BudgetCategoryOptions struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}