package services

import (
	"strconv"
	"wow-bato-backend/internal/models"

	"gorm.io/gorm"
)

type BudgetItemService struct {
	db *gorm.DB
}

func NewBudgetItemService (db *gorm.DB) *BudgetItemService {
	return &BudgetItemService{db: db}
}

func (s *BudgetItemService) AddBudgetItem(projectID string, budgetItem models.NewBudgetItem) error {

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

	result := s.db.Create(&newBudgetItem)

	return result.Error
}

func (s *BudgetItemService) GetAllBudgetItem(projectID string, filter string, page string) ([]models.Budget_Item, error) {

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
		if err := s.db.Where("project_id = ?", projectID_int).Find(&budgetItem).Limit(limit).Offset(offset).Error; err != nil {
			return []models.Budget_Item{}, err
		}

		return budgetItem, nil
	}

	if err := s.db.Where("project_id = ? AND status = ?", projectID_int, filter).Find(&budgetItem).Limit(limit).Offset(offset).Error; err != nil {
		return []models.Budget_Item{}, err
	}

	return budgetItem, nil
}

func (s *BudgetItemService) CountBudgetItem(projectID string) (int64, error) {

	projectID_int, err := strconv.Atoi(projectID)
	if err != nil {
		return 0, err
	}

	var count int64
	if err := s.db.Model(&models.Budget_Item{}).Where("project_id = ?", projectID_int).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (s *BudgetItemService) GetSingleBudgetItem(categoryID string, budgetItemID string) (models.Budget_Item, error) {

	categoryID_int, err := strconv.Atoi(categoryID)
	if err != nil {
		return models.Budget_Item{}, err
	}

	budgetItemID_int, err := strconv.Atoi(budgetItemID)
	if err != nil {
		return models.Budget_Item{}, err
	}

	var budgetItem models.Budget_Item
	if err := s.db.Where("categoryID = ? AND status = ?", categoryID_int, budgetItemID_int).First(&budgetItem).Error; err != nil {
		return models.Budget_Item{}, err
	}

	return budgetItem, nil
}

func (s *BudgetItemService) UpdateBudgetItemStatus(budgetItemID string, newStatus models.UpdateStatus) error {

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
	if err := s.db.Where("id = ?", budgetItemID_int).First(&budgetItem).Error; err != nil {
		return err
	}

	budgetItem.Status = updateStatus

	result := s.db.Save(&budgetItem)

	return result.Error
}

func (s *BudgetItemService) DeleteBudgetItem(budgetItemID string) error {

	budgetItemID_int, err := strconv.Atoi(budgetItemID)
	if err != nil {
		return err
	}

	if err := s.db.Where("id = ?", budgetItemID_int).Delete(&models.Budget_Item{}).Error; err != nil {
		return err
	}

	return nil
}
