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
	"errors"
	"fmt"
	"strconv"
	database "wow-bato-backend/internal"
	"wow-bato-backend/internal/models"
)

// Domain-specific errors for budget category operations
var (
	ErrBudgetCategoryNotFound            = errors.New("budget category not found")
	ErrInvalidBudgetCategoryID           = errors.New("invalid budget category ID format")
	ErrEmptyBudgetCategoryName           = errors.New("budget category name cannot be empty")
	ErrEmptyBudgetDescription            = errors.New("budget category description cannot be empty")
	ErrInvalidBudgetBarangayID           = errors.New("invalid barangay ID for budget category")
	ErrorInvalidBarangayIDBudgetCategory = errors.New("invalid barangay ID format")
)

// validateBudgetCategory validates required fields for budget category operations
func validateBudgetCategory(category models.NewBudgetCategory) error {
	if category.Name == "" {
		return ErrEmptyBudgetCategoryName
	}
	if category.Description == "" {
		return ErrEmptyBudgetDescription
	}
	if category.Barangay_ID == 0 {
		return ErrInvalidBudgetBarangayID
	}
	return nil
}

// AddBudgetCategory creates a new budget category record with validated information.
//
// This function implements the core business logic for budget category registration,
// ensuring data consistency and proper validation before storage.
//
// Parameters:
//   - budgetCategory: models.NewBudgetCategory - The budget category data containing:
//   - Name: Name of the budget category (required)
//   - Description: Detailed description of the budget category (required)
//   - Barangay_ID: ID of the associated barangay (required)
//
// Returns:
//   - error: nil on successful creation, or an error describing the failure:
//   - ErrDatabaseConnection: When database connection fails
//   - ErrValidation: When required fields are missing or invalid
//   - ErrDatabaseOperation: When budget category creation fails
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
	// Validate input data
	if err := validateBudgetCategory(budgetCategory); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	db, err := database.ConnectDB()
	if err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}

	newBudgetCategory := models.Budget_Category{
		Name:        budgetCategory.Name,
		Description: budgetCategory.Description,
		Barangay_ID: budgetCategory.Barangay_ID,
	}

	if result := db.Create(&newBudgetCategory); result.Error != nil {
		return fmt.Errorf("failed to create budget category: %w", result.Error)
	}

	return nil
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
//   - ErrInvalidID: When budget_ID format is invalid
//   - ErrNotFound: When budget category does not exist
//   - ErrDependencyExists: When budget category has associated records
//   - ErrDatabaseOperation: When deletion operation fails
//
// Example usage:
//
//	if err := DeleteBudgetCategory("123"); err != nil {
//	    return fmt.Errorf("failed to delete budget category: %w", err)
//	}
func DeleteBudgetCategory(budget_ID string) error {
	db, err := database.ConnectDB()
	if err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}

	budget_ID_int, err := strconv.Atoi(budget_ID)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidBudgetCategoryID, budget_ID)
	}

	result := db.Where("id = ?", budget_ID_int).Delete(&models.Budget_Category{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete budget category: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("%w: ID %d", ErrBudgetCategoryNotFound, budget_ID_int)
	}

	return nil
}

// UpdateBudgetCategory modifies existing budget category information.
//
// This function implements validation and update logic for budget category records,
// ensuring data consistency and proper error handling.
//
// Parameters:
//   - budget_ID: string - Unique identifier of the budget category to update
//   - updateBudgetCategory: models.UpdateBudgetCategory - The update data containing:
//   - Name: Updated name of the budget category
//   - Description: Updated description of the budget category
//
// Returns:
//   - error: nil on successful update, or an error describing the failure:
//   - ErrInvalidID: When budget_ID format is invalid
//   - ErrNotFound: When budget category does not exist
//   - ErrValidation: When update data is invalid
//   - ErrDatabaseOperation: When update operation fails
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
	// Validate update data
	if updateBudgetCategory.Name == "" {
		return ErrEmptyBudgetCategoryName
	}
	if updateBudgetCategory.Description == "" {
		return ErrEmptyBudgetDescription
	}

	db, err := database.ConnectDB()
	if err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}

	budget_ID_int, err := strconv.Atoi(budget_ID)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidBudgetCategoryID, budget_ID)
	}

	var budgetCategory models.Budget_Category
	if err := db.Where("id = ?", budget_ID_int).First(&budgetCategory).Error; err != nil {
		return fmt.Errorf("%w: ID %d", ErrBudgetCategoryNotFound, budget_ID_int)
	}

	budgetCategory.Name = updateBudgetCategory.Name
	budgetCategory.Description = updateBudgetCategory.Description

	if result := db.Save(&budgetCategory); result.Error != nil {
		return fmt.Errorf("failed to update budget category: %w", result.Error)
	}

	return nil
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
//   - ID: Unique identifier of the budget category
//   - Name: Name of the budget category
//   - Description: Detailed description
//   - CreatedAt: Timestamp of record creation
//   - UpdatedAt: Timestamp of last update
//   - error: nil on successful retrieval, or an error describing the failure:
//   - ErrInvalidID: When barangay_ID format is invalid
//   - ErrInvalidPagination: When limit or page parameters are invalid
//   - ErrDatabaseOperation: When retrieval operation fails
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
		return nil, fmt.Errorf("database connection failed: %w", err)
	}

	barangay_ID_int, err := strconv.Atoi(barangay_ID)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrorInvalidBarangayIDBudgetCategory, barangay_ID)
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return nil, fmt.Errorf("invalid limit parameter: %w", err)
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return nil, fmt.Errorf("invalid page parameter: %w", err)
	}

	if pageInt < 1 {
		return nil, fmt.Errorf("page number must be greater than zero")
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

	if result.Error != nil {
		return nil, fmt.Errorf("failed to retrieve budget categories: %w", result.Error)
	}

	return budgetCategory, nil
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
//   - ErrInvalidID: When barangay_ID format is invalid
//   - ErrNotFound: When barangay does not exist
//   - ErrDatabaseOperation: When count operation fails
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
		return 0, fmt.Errorf("database connection failed: %w", err)
	}

	barangay_ID_int, err := strconv.Atoi(barangay_ID)
	if err != nil {
		return 0, fmt.Errorf("%w: %s", ErrorInvalidBarangayIDBudgetCategory, barangay_ID)
	}

	var count int64
	if err := db.Model(&models.Budget_Category{}).Where("barangay_ID = ?", barangay_ID_int).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count budget categories: %w", err)
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
//   - ID: Unique identifier of the budget category
//   - Name: Name of the budget category
//   - Description: Detailed description
//   - CreatedAt: Timestamp of record creation
//   - UpdatedAt: Timestamp of last update
//   - error: nil on successful retrieval, or an error describing the failure:
//   - ErrInvalidID: When budget_ID format is invalid
//   - ErrNotFound: When budget category does not exist
//   - ErrDatabaseOperation: When retrieval operation fails
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
		return models.DisplayBudgetCategory{}, fmt.Errorf("database connection failed: %w", err)
	}

	budget_ID_int, err := strconv.Atoi(budget_ID)
	if err != nil {
		return models.DisplayBudgetCategory{}, fmt.Errorf("%w: %s", ErrInvalidBudgetCategoryID, budget_ID)
	}

	var budgetCategory models.DisplayBudgetCategory
	result := db.Model(&models.Budget_Category{}).
		Where("barangay_ID = ? AND id = ?", barangay_ID, budget_ID_int).
		Scan(&budgetCategory)

	if result.Error != nil {
		return models.DisplayBudgetCategory{}, fmt.Errorf("failed to retrieve budget category: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return models.DisplayBudgetCategory{}, fmt.Errorf("%w: ID %d", ErrBudgetCategoryNotFound, budget_ID_int)
	}

	return budgetCategory, nil
}
