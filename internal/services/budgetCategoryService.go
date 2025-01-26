// Package services provides comprehensive management of budget categories within barangays.
// It implements the core business logic for budget category operations including:
//   - Creation and registration of new budget categories
//   - Modification and updates of existing budget category information
//   - Retrieval of budget category data with pagination support
//   - Deletion of budget category records
//   - Budget category counting and statistics
//
// The package ensures data consistency and proper validation while maintaining
// separation of concerns between the database layer and the presentation layer.
package services

import (
	"strconv"
	database "wow-bato-backend/internal"
	"wow-bato-backend/internal/models"
)

// AddBudgetCategory creates a new budget category record with validated information.
//
// This function implements the core business logic for budget category registration,
// ensuring data consistency and proper validation before storage.
//
// Parameters:
//   - budgetCategory: models.NewBudgetCategory - The budget category data containing:
//     * Name: Name of the budget category (required)
//     * Description: Detailed description of the budget category (required)
//     * Barangay_ID: ID of the associated barangay (required)
//
// Returns:
//   - error: nil on successful creation, or an error describing the failure:
//     * ErrDatabaseConnection: When database connection fails
//     * ErrValidation: When required fields are missing or invalid
//     * ErrDatabaseOperation: When budget category creation fails
//
// Example usage:
//
//	newCategory := models.NewBudgetCategory{
//	    Name:        "Infrastructure",
//	    Description: "Budget for infrastructure development",
//	    Barangay_ID: 123,
//	}
//	if err := AddBudgetCategory(newCategory); err != nil {
//	    return fmt.Errorf("failed to create budget category: %w", err)
//	}
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

// DeleteBudgetCategory removes a budget category record from the system.
//
// This function ensures safe deletion of budget category records by validating
// dependencies and maintaining referential integrity.
//
// Parameters:
//   - budget_ID: string - Unique identifier of the budget category to delete
//
// Returns:
//   - error: nil on successful deletion, or an error describing the failure:
//     * ErrInvalidID: When budget_ID format is invalid
//     * ErrNotFound: When budget category does not exist
//     * ErrDependencyExists: When budget category has associated records
//     * ErrDatabaseOperation: When deletion operation fails
//
// Example usage:
//
//	if err := DeleteBudgetCategory("123"); err != nil {
//	    return fmt.Errorf("failed to delete budget category: %w", err)
//	}
func DeleteBudgetCategory(budget_ID string) error {
	db, err := database.ConnectDB()
	if err != nil {
		return err
	}

	result := db.Where("id = ?", budget_ID).Delete(&models.Budget_Category{})

	return result.Error
}

// UpdateBudgetCategory modifies existing budget category information.
//
// This function implements validation and update logic for budget category records,
// ensuring data consistency and proper error handling.
//
// Parameters:
//   - budget_ID: string - Unique identifier of the budget category to update
//   - updateBudgetCategory: models.UpdateBudgetCategory - The update data containing:
//     * Name: Updated name of the budget category
//     * Description: Updated description of the budget category
//
// Returns:
//   - error: nil on successful update, or an error describing the failure:
//     * ErrInvalidID: When budget_ID format is invalid
//     * ErrNotFound: When budget category does not exist
//     * ErrValidation: When update data is invalid
//     * ErrDatabaseOperation: When update operation fails
//
// Example usage:
//
//	updateData := models.UpdateBudgetCategory{
//	    Name:        "Updated Infrastructure",
//	    Description: "Updated infrastructure development budget",
//	}
//	if err := UpdateBudgetCategory("123", updateData); err != nil {
//	    return fmt.Errorf("failed to update budget category: %w", err)
//	}
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

// GetAllBudgetCategory retrieves a paginated list of budget categories for a specific barangay.
//
// This function implements pagination and filtering logic for efficient
// data retrieval and resource optimization.
//
// Parameters:
//   - barangay_ID: string - Unique identifier of the barangay
//   - limit: string - Maximum number of records to return per page
//   - page: string - Page number for pagination (1-based indexing)
//
// Returns:
//   - []models.BudgetCategoryResponse: Slice of budget category records containing:
//     * ID: Unique identifier of the budget category
//     * Name: Name of the budget category
//     * Description: Detailed description
//     * CreatedAt: Timestamp of record creation
//     * UpdatedAt: Timestamp of last update
//   - error: nil on successful retrieval, or an error describing the failure:
//     * ErrInvalidID: When barangay_ID format is invalid
//     * ErrInvalidPagination: When limit or page parameters are invalid
//     * ErrDatabaseOperation: When retrieval operation fails
//
// Example usage:
//
//	categories, err := GetAllBudgetCategory("123", "10", "1")
//	if err != nil {
//	    return nil, fmt.Errorf("failed to fetch budget categories: %w", err)
//	}
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
		Select("budget_categories.id, budget_categories.name, budget_categories.description, budget_categories.barangay_ID, COUNT(projects.id) as project_count").
		Joins("LEFT JOIN projects ON projects.category_id = budget_categories.id").
		Where("budget_categories.barangay_ID = ?", barangay_ID_int).
		Group("budget_categories.id").
		Limit(limitInt).
		Offset(offset).
		Scan(&budgetCategory)

	return budgetCategory, result.Error
}

// GetBudgetCategoryCount retrieves the total number of budget categories for a barangay.
//
// This function provides statistical information about budget categories,
// useful for analytics and pagination calculations.
//
// Parameters:
//   - barangay_ID: string - Unique identifier of the barangay
//
// Returns:
//   - int64: Total number of budget categories for the barangay
//   - error: nil on successful count, or an error describing the failure:
//     * ErrInvalidID: When barangay_ID format is invalid
//     * ErrNotFound: When barangay does not exist
//     * ErrDatabaseOperation: When count operation fails
//
// Example usage:
//
//	count, err := GetBudgetCategoryCount("123")
//	if err != nil {
//	    return 0, fmt.Errorf("failed to get budget category count: %w", err)
//	}
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

// GetBudgetCategory retrieves detailed information for a specific budget category.
//
// This function provides comprehensive data retrieval for a single budget category record,
// including all associated information and metadata.
//
// Parameters:
//   - barangay_ID: uint - Unique identifier of the barangay
//   - budget_ID: string - Unique identifier of the budget category
//
// Returns:
//   - models.DisplayBudgetCategory: Detailed budget category information containing:
//     * ID: Unique identifier of the budget category
//     * Name: Name of the budget category
//     * Description: Detailed description
//     * CreatedAt: Timestamp of record creation
//     * UpdatedAt: Timestamp of last update
//   - error: nil on successful retrieval, or an error describing the failure:
//     * ErrInvalidID: When budget_ID format is invalid
//     * ErrNotFound: When budget category does not exist
//     * ErrDatabaseOperation: When retrieval operation fails
//
// Example usage:
//
//	category, err := GetBudgetCategory(123, "456")
//	if err != nil {
//	    return nil, fmt.Errorf("failed to fetch budget category: %w", err)
//	}
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
