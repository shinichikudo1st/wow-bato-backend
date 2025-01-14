// Package services provides budget item-related business logic and operations for the application.
// It handles budget item management while maintaining separation of concerns from the database and presentation layers.
package services

import (
	"strconv"
	database "wow-bato-backend/internal"
	"wow-bato-backend/internal/models"
)

// AddBudgetItem adds a new budget item to the database
//
// This function performs the following operations:
// 1. Establishes a database connection
// 2. Converts the categoryID from string to int
// 3. Creates a new budget item record in the database
// 4. Returns nil if successful, otherwise returns an error
//
// Parameters:
//   - categoryID: string - The category ID of the budget item
//   - budgetItem: models.NewBudgetItem - The new budget item data
//
// Returns:
//   - error: Returns nil if successful, otherwise returns an error:
//   - Database connection errors
//   - Database creation errors
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
		Name:             budgetItem.Name,
		Amount_Allocated: budgetItem.Amount_Allocated,
		Description:      budgetItem.Description,
		Status:           budgetItem.Status,
		CategoryID:       uint(categoryID_int),
	}

	result := db.Create(&newBudgetItem)

	return result.Error
}

// GetAllBudgetItem retrieves all budget items for a project
//
// This function performs the following operations:
// 1. Establishes a database connection
// 2. Converts the projectID from string to int
// 3. Retrieves all budget items for the project from the database
// 4. Returns the retrieved budget items and nil if successful, otherwise returns an error
//
// Parameters:
//   - projectID: string - The ID of the project
//   - filter: string - The status of the budget item
//
// Returns:
//   - []models.Budget_Item: Returns the retrieved budget items
//   - error: Returns nil if successful, otherwise returns an error
func GetAllBudgetItem(projectID string, filter string) ([]models.Budget_Item, error) {
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

// GetSingleBudgetItem retrieves a single budget item from the database
//
// This function performs the following operations:
// 1. Establishes a database connection
// 2. Converts the categoryID and budgetItemID from string to int
// 3. Retrieves the budget item record from the database that matches the categoryID and budgetItemID
// 4. Returns the retrieved budget item and nil if successful, otherwise returns an error
//
// Parameters:
//   - categoryID: string - The category ID of the budget item
//   - budgetItemID: string - The ID of the budget item
//
// Returns:
//   - models.Budget_Item: Returns the retrieved budget item
//   - error: Returns nil if successful, otherwise returns an error
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

// UpdateBudgetItemStatus updates the status of a budget item in the database
//
// This function performs the following operations:
// 1. Establishes a database connection
// 2. Converts the budgetItemID from string to int
// 3. Retrieves the budget item record from the database that matches the budgetItemID
// 4. Updates the status of the budget item record in the database
// 5. Returns nil if successful, otherwise returns an error
//
// Parameters:
//   - budgetItemID: string - The ID of the budget item
//   - newStatus: models.UpdateStatus - The new status of the budget item
//
// Returns:
//   - error: Returns nil if successful, otherwise returns an error
//   - Database connection errors
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
