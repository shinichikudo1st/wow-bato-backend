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

// For displaying at the client side BudgetCategoryList.tsx
type BudgetCategoryResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Barangay_ID uint   `json:"barangay_ID"`
	ProjectCount int64  `json:"project_count"`
}

// Display projects inside BudgetCategoryList
type ProjectResponse struct {
	Name        string `json:"name"`
	Status      string `json:"status"`
	Category_ID uint   `json:"category_ID"`
}

type DisplayBudgetCategory struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
