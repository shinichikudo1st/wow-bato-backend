package services

import (
	database "wow-bato-backend/internal"
	"wow-bato-backend/internal/models"
)

func AddBudgetCategory(budgetCategory models.NewBudgetCategory) error {
	db, err := database.ConnectDB()
	if err != nil {
		return err
	}

	newBudgetCategory := models.Budget_Category{
		Name: budgetCategory.Name,
		Description: budgetCategory.Description,
		Barangay_ID: budgetCategory.Barangay_ID,
	}

	result := db.Create(&newBudgetCategory)

	return result.Error
}

func DeleteBudgetCategory(budgetID string) error {
	db, err := database.ConnectDB()
	if err != nil {
		return err
	}

	result := db.Where("id = ?", budgetID).Delete(&models.Budget_Category{})

	return result.Error
}