package services

import (
	"errors"
	"fmt"
	"strconv"
	database "wow-bato-backend/internal"
	"wow-bato-backend/internal/models"
)

var (
	ErrBudgetCategoryNotFound            = errors.New("budget category not found")
	ErrInvalidBudgetCategoryID           = errors.New("invalid budget category ID format")
	ErrEmptyBudgetCategoryName           = errors.New("budget category name cannot be empty")
	ErrEmptyBudgetDescription            = errors.New("budget category description cannot be empty")
	ErrInvalidBudgetBarangayID           = errors.New("invalid barangay ID for budget category")
	ErrorInvalidBarangayIDBudgetCategory = errors.New("invalid barangay ID format")
)

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

func AddBudgetCategory(budgetCategory models.NewBudgetCategory) error {
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

func UpdateBudgetCategory(budget_ID string, updateBudgetCategory models.UpdateBudgetCategory) error {
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
