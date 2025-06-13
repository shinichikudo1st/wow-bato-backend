package services

import (
	"wow-bato-backend/internal/models"

	"gorm.io/gorm"
)

type PublicDashboardStats struct {
	TotalProjects    int64
	TotalUsers       int64
	TotalBudgetItems int64
	ProjectsByStatus map[string]int64
	UsersByRole      map[string]int64
	BudgetByProject  map[string]float64
}

type PublicDashboardService struct {
	db *gorm.DB
}

func NewPublicDashboardService(db *gorm.DB) *PublicDashboardService {
	return &PublicDashboardService{db: db}
}

func (s *PublicDashboardService) GetPublicDashboardStats() (PublicDashboardStats, error) {
	var stats PublicDashboardStats
	stats.ProjectsByStatus = make(map[string]int64)
	stats.UsersByRole = make(map[string]int64)
	stats.BudgetByProject = make(map[string]float64)

	// Total Projects
	s.db.Model(&models.Project{}).Count(&stats.TotalProjects)
	// Total Users
	s.db.Model(&models.User{}).Count(&stats.TotalUsers)
	// Total Budget Items
	s.db.Model(&models.Budget_Item{}).Count(&stats.TotalBudgetItems)

	// Projects by Status
	var statusResults []struct {
		Status string
		Count  int64
	}
	s.db.Model(&models.Project{}).Select("status, COUNT(*) as count").Group("status").Scan(&statusResults)
	for _, r := range statusResults {
		stats.ProjectsByStatus[r.Status] = r.Count
	}

	// Users by Role
	var roleResults []struct {
		Role  string
		Count int64
	}
	s.db.Model(&models.User{}).Select("role, COUNT(*) as count").Group("role").Scan(&roleResults)
	for _, r := range roleResults {
		stats.UsersByRole[r.Role] = r.Count
	}

	// Budget Allocation by Project
	var budgetResults []struct {
		ProjectName string
		TotalBudget float64
	}
	s.db.Table("budget_items").
		Select("projects.name as project_name, SUM(budget_items.amount_allocated) as total_budget").
		Joins("JOIN projects ON projects.id = budget_items.project_id").
		Group("projects.name").
		Scan(&budgetResults)
	for _, r := range budgetResults {
		stats.BudgetByProject[r.ProjectName] = r.TotalBudget
	}

	return stats, nil
}
