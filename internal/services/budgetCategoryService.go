package services

import (
	"strconv"
	database "wow-bato-backend/internal"
	"wow-bato-backend/internal/models"

	"gorm.io/gorm"
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

func GetAllBudgetCategory(barangay_ID string, limit string, page string) ([]models.Budget_Category, error) {
	db, err := database.ConnectDB()
	if err != nil {
		return []models.Budget_Category{}, err
	}

	barangay_ID_int, err := strconv.Atoi(barangay_ID)
	if err != nil {
		return []models.Budget_Category{}, err
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return []models.Budget_Category{}, err
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return []models.Budget_Category{}, err
	}

	offset := (pageInt - 1) * limitInt

	var budgetCategory []models.Budget_Category
	result := db.Model(&models.Budget_Category{}).
		Preload("Projects", func(db *gorm.DB) *gorm.DB{
			return db.Select("name, status, category_id")
		}).
		Select("id, name, description, barangay_ID").
		Where("barangay_ID = ?", barangay_ID_int).
		Limit(limitInt).
		Offset(offset).
		Find(&budgetCategory)


	return budgetCategory, result.Error 
}


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

func GetBudgetCategory(barangay_ID uint, budget_ID string)(models.DisplayBudgetCategory, error){
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