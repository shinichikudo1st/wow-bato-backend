package services

import (
	"strconv"
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

func DeleteBudgetCategory(budget_ID string) error {
	db, err := database.ConnectDB()
	if err != nil {
		return err
	}

	result := db.Where("id = ?", budget_ID).Delete(&models.Budget_Category{})

	return result.Error
}

func UpdateBudgetCategory(budget_ID string, updateBudgetCategory models.UpdateBudgetCategory) error {
	db, err := database.ConnectDB()
	if err != nil {
		return err
	}

	budget_ID_int, err := strconv.Atoi(budget_ID)
	if err != nil {
		return err
	}

	var budgetCategory models.Budget_Category
	if err := db.Where("id = ?", budget_ID_int).First(&budgetCategory).Error; err != nil {
		return err
	}

	budgetCategory.Name = updateBudgetCategory.Name
	budgetCategory.Description = updateBudgetCategory.Description

	result := db.Save(&budgetCategory)

	return result.Error

}

func GetAllBudgetCategory(barangay_ID string, limit string, page string) ([]models.BudgetCategoryResponse, error) {
	db, err := database.ConnectDB()
	if err != nil {
		return []models.BudgetCategoryResponse{}, err
	}

	barangay_ID_int, err := strconv.Atoi(barangay_ID)
	if err != nil {
		return []models.BudgetCategoryResponse{}, err
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return []models.BudgetCategoryResponse{}, err
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return []models.BudgetCategoryResponse{}, err
	}

	offset := (pageInt - 1) * limitInt

	var budgetCategory []models.BudgetCategoryResponse
	result := db.Model(&models.Budget_Category{}).
		Select("id, name, description, barangay_ID").
		Where("barangay_ID = ?", barangay_ID_int).
		Limit(limitInt).
		Offset(offset).
		Scan(&budgetCategory)


	return budgetCategory, result.Error 
}

func GetBudgetCategory(barangay_ID string, budget_ID string)(models.Budget_Category, error){
	db, err := database.ConnectDB()
	if err != nil {
		return models.Budget_Category{}, err
	}

	barangay_ID_int, err := strconv.Atoi(barangay_ID)
	if err != nil {
		return models.Budget_Category{}, err
	}

	budget_ID_int, err := strconv.Atoi(budget_ID)
	if err != nil {
		return models.Budget_Category{}, err
	}

	var budgetCategory models.Budget_Category
	result := db.Where("barangay_ID = ? AND id = ?", barangay_ID_int, budget_ID_int).First(&budgetCategory)
	

	return budgetCategory, result.Error
}