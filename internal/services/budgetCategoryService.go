// Package services provides budget category-related business logic and operations for the application.
// It handles budget category management while maintaining separation of concerns from the database and presentation layers.
package services

import (
	"strconv"
	database "wow-bato-backend/internal"
	"wow-bato-backend/internal/models"
)

// AddBudgetCategory adds a new budget category
//
// This function performs the following operations:
// 1. Establishes a database connection
// 2. Creates a new budget category record in the database
// 3. Returns nil if successful, otherwise returns an error
//
// Parameters:
//   - budgetCategory: models.NewBudgetCategory -
//     Contains budget category data including:
//   - Name: Budget category namse
//   - Description: Budget category description
//   - Barangay_ID: Barangay ID
//
// Returns:
//   - error: Returns nil if successful, otherwise returns an error:
//   - Database connection errors
//   - Database creation errors
func AddBudgetCategory(budgetCategory models.NewBudgetCategory) error {
	db, err := database.ConnectDB()
	if err != nil {
		return err
	}

	newBudgetCategory := models.Budget_Category{
		Name:        budgetCategory.Name,
		Description: budgetCategory.Description,
		Barangay_ID: budgetCategory.Barangay_ID,
	}

	result := db.Create(&newBudgetCategory)

	return result.Error
}

// DeleteBudgetCategory deletes a budget category
//
// This function performs the following operations:
// 1. Establishes a database connection
// 2. Deletes a budget category record from the database
// 3. Returns nil if successful, otherwise returns an error
//
// Parameters:
//   - budget_ID: string - Budget category ID
//
// Returns:
//   - error: Returns nil if successful, otherwise returns an error:
//   - Database connection errors
//   - Database deletion errors
func DeleteBudgetCategory(budget_ID string) error {
	db, err := database.ConnectDB()
	if err != nil {
		return err
	}

	result := db.Where("id = ?", budget_ID).Delete(&models.Budget_Category{})

	return result.Error
}

// UpdateBudgetCategory updates a budget category
//
// This function performs the following operations:
// 1. Establishes a database connection
// 2. Converts the budget_ID from string to int
// 3. Updates a budget category record in the database that matches the budget_ID
// 4. Returns nil if successful, otherwise returns an error
//
// Parameters:
//   - budget_ID: string - Budget category ID
//   - updateBudgetCategory: models.UpdateBudgetCategory -
//     Contains budget category data including:
//   - Name: Budget category name
//   - Description: Budget category description
//
// Returns:
//   - error: Returns nil if successful, otherwise returns an error:
//   - Database connection errors
//   - Database update errors
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

// GetAllBudgetCategory retrieves all budget categories from the database
//
// This function performs the following operations:
// 1. Establishes a database connection
// 2. Converts the limit and page from string to int
// 3. Retrieves all budget categories from the database
// 4. Returns the budget categories and any errors
//
// Parameters:
//   - limit: string - The number of budget categories to retrieve per page
//   - page: string - The page number to retrieve
//
// Returns:
//   - []models.BudgetCategoryResponse: A slice of budget category responses
//   - error: Returns nil if successful, otherwise returns an error:
//   - Database connection errors
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

// GetBudgetCategoryCount returns the number of budget categories for a barangay
//
// This function performs the following operations:
// 1. Establishes a database connection
// 2. Converts the barangay_ID from string to int
// 3. Retrieves the number of budget categories for the barangay
// 4. Returns the number of budget categories and any errors
//
// Parameters:
//   - barangay_ID: string - The ID of the barangay
//
// Returns:
//   - int64: The number of budget categories
//   - error: Returns nil if successful, otherwise returns an error:
//   - Database connection errors
func GetBudgetCategoryCount(barangay_ID string) (int64, error) {
	db, err := database.ConnectDB()
	if err != nil {
		return 0, err
	}

	barangay_ID_int, err := strconv.Atoi(barangay_ID)
	if err != nil {
		return 0, err
	}

	var count int64
	if err := db.Model(&models.Budget_Category{}).Where("barangay_ID = ?", barangay_ID_int).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// GetBudgetCategory returns a single budget category
//
// This function performs the following operations:
// 1. Establishes a database connection
// 2. Converts the budget_ID from string to int
// 3. Retrieves the budget category record from the database that matches the budget_ID
// 4. Returns the budget category and any errors
//
// Parameters:
//   - barangay_ID: uint - The ID of the barangay
//   - budget_ID: string - The ID of the budget category
//
// Returns:
//   - models.DisplayBudgetCategory: A budget category response
//   - error: Returns nil if successful, otherwise returns an error:
//   - Database connection errors
func GetBudgetCategory(barangay_ID uint, budget_ID string) (models.DisplayBudgetCategory, error) {
	db, err := database.ConnectDB()
	if err != nil {
		return models.DisplayBudgetCategory{}, err
	}

	budget_ID_int, err := strconv.Atoi(budget_ID)
	if err != nil {
		return models.DisplayBudgetCategory{}, err
	}

	var budgetCategory models.DisplayBudgetCategory
	result := db.Model(&models.Budget_Category{}).Where("barangay_ID = ? AND id = ?", barangay_ID, budget_ID_int).Scan(&budgetCategory)

	return budgetCategory, result.Error
}
