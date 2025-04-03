package services

import (
	"strconv"
	database "wow-bato-backend/internal"
	"wow-bato-backend/internal/models"
)

func AddBudgetItem(projectID string, budgetItem models.NewBudgetItem) error {
	db, err := database.ConnectDB()
	if err != nil {
		return err
	}

	projectID_int, err := strconv.Atoi(projectID)
	if err != nil {
		return err
	}

	newBudgetItem := models.Budget_Item{
		Name:             budgetItem.Name,
		Amount_Allocated: budgetItem.Amount_Allocated,
		Description:      budgetItem.Description,
		Status:           budgetItem.Status,
		ProjectID:       uint(projectID_int),
	}

	result := db.Create(&newBudgetItem)

	return result.Error
}

func GetAllBudgetItem(projectID string, filter string, page string) ([]models.Budget_Item, error) {
	db, err := database.ConnectDB()
	if err != nil {
		return []models.Budget_Item{}, err
	}

	projectID_int, err := strconv.Atoi(projectID)
	if err != nil {
		return []models.Budget_Item{}, err
	}

	page_int, err := strconv.Atoi(page)
	if err != nil {
		return []models.Budget_Item{}, err
	}

	limit := 5 // temporary hardcoded

	offset := (page_int - 1) * limit


	var budgetItem []models.Budget_Item
	if filter == "All" {
		if err := db.Where("project_id = ?", projectID_int).Find(&budgetItem).Limit(limit).Offset(offset).Error; err != nil {
			return []models.Budget_Item{}, err
		}

		return budgetItem, nil
	}

	if err := db.Where("project_id = ? AND status = ?", projectID_int, filter).Find(&budgetItem).Limit(limit).Offset(offset).Error; err != nil {
		return []models.Budget_Item{}, err
	}

	return budgetItem, nil
}

func CountBudgetItem(projectID string) (int64, error) {
	db, err := database.ConnectDB()
	if err != nil {
		return 0, err
	}

	projectID_int, err := strconv.Atoi(projectID)
	if err != nil {
		return 0, err
	}

	var count int64
	if err := db.Model(&models.Budget_Item{}).Where("project_id = ?", projectID_int).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func GetSingleBudgetItem(categoryID string, budgetItemID string) (models.Budget_Item, error) {
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

	var updateStatus string
	if newStatus.Status == "approve" {
		updateStatus = "Approved"
	} else if newStatus.Status == "reject" {
		updateStatus = "Rejected"
	}

	budgetItemID_int, err := strconv.Atoi(budgetItemID)
	if err != nil {
		return err
	}

	var budgetItem models.Budget_Item
	if err := db.Where("id = ?", budgetItemID_int).First(&budgetItem).Error; err != nil {
		return err
	}

	budgetItem.Status = updateStatus

	result := db.Save(&budgetItem)

	return result.Error
}

func DeleteBudgetItem(budgetItemID string) error {
	db, err := database.ConnectDB()
	if err != nil {
		return err
	}

	budgetItemID_int, err := strconv.Atoi(budgetItemID)
	if err != nil {
		return err
	}

	if err := db.Where("id = ?", budgetItemID_int).Delete(&models.Budget_Item{}).Error; err != nil {
		return err
	}

	return nil
}
