package services

import (
	"strconv"
	database "wow-bato-backend/internal"
	"wow-bato-backend/internal/models"
)

func AddBudgetItem(categoryID string, budgetItem models.NewBudgetItem) error {
	db, err := database.ConnectDB()
	if err != nil {
		return err
	}

	categoryID_int, err := strconv.Atoi(categoryID)
	if err != nil {
		return err
	}

	newBudgetItem := models.Budget_Item{
		Name: budgetItem.Name,
		Amount_Allocated: budgetItem.Amount_Allocated,
		Amount_Spent: budgetItem.Amount_Spent,
		Description: budgetItem.Description,
		Status: budgetItem.Status,
		CategoryID: uint(categoryID_int),
	}

	result :=  db.Create(&newBudgetItem)

	return result.Error
}