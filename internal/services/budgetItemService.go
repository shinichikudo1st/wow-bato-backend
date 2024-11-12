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
		Description: budgetItem.Description,
		Status: budgetItem.Status,
		CategoryID: uint(categoryID_int),
	}

	result :=  db.Create(&newBudgetItem)

	return result.Error
}

func GetAllBudgetItem(projectID string, filter string) ([]models.Budget_Item, error){
	db, err := database.ConnectDB()
	if err != nil {
		return []models.Budget_Item{}, err
	}

	projectID_int, err := strconv.Atoi(projectID)
	if err != nil {
		return []models.Budget_Item{}, err
	}

	var budgetItem []models.Budget_Item
	if err := db.Where("ProjectID = ? AND status = ?", projectID_int, filter).Find(&budgetItem).Error; err != nil {
		return []models.Budget_Item{}, err
	}

	return budgetItem, nil
}

func GetSingleBudgetItem(categoryID string, budgetItemID string)(models.Budget_Item, error){
	db, err := database.ConnectDB()
	if err != nil {
		return models.Budget_Item{}, err
	}

	categoryID_int, err := strconv.Atoi(categoryID)
	if err != nil {
		return models.Budget_Item{}, err
	}

	budgetItemID_int, err := strconv.Atoi(budgetItemID)
	if err != nil {
		return models.Budget_Item{}, err
	}

	var budgetItem models.Budget_Item
	if err := db.Where("categoryID = ? AND status = ?", categoryID_int, budgetItemID_int).First(&budgetItem).Error; err != nil {
		return models.Budget_Item{}, err
	}

	return budgetItem, nil
}

func UpdateBudgetItemStatus(budgetItemID string, newStatus models.UpdateStatus) error {
	db, err := database.ConnectDB()
	if err != nil {
		return err
	}

	budgetItemID_int, err := strconv.Atoi(budgetItemID)
	if err != nil {
		return err
	}

	var budgetItem models.Budget_Item
	if err := db.Where("id = ?", budgetItemID_int).First(&budgetItem).Error; err != nil {
		return err
	}

	budgetItem.Status = newStatus.Status

	result := db.Save(&budgetItem)

	return result.Error
}