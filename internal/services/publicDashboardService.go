package services

import (
	"time"
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

	s.db.Model(&models.Project{}).Count(&stats.TotalProjects)
	s.db.Model(&models.User{}).Count(&stats.TotalUsers)
	s.db.Model(&models.Budget_Item{}).Count(&stats.TotalBudgetItems)

	var statusResults []struct {
		Status string
		Count  int64
	}
	s.db.Model(&models.Project{}).Select("status, COUNT(*) as count").Group("status").Scan(&statusResults)
	for _, r := range statusResults {
		stats.ProjectsByStatus[r.Status] = r.Count
	}

	var roleResults []struct {
		Role  string
		Count int64
	}
	s.db.Model(&models.User{}).Select("role, COUNT(*) as count").Group("role").Scan(&roleResults)
	for _, r := range roleResults {
		stats.UsersByRole[r.Role] = r.Count
	}

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

type CompleteStats struct {
	Complete   int64 `json:"complete"`
	Incomplete int64 `json:"incomplete"`
}

func (s *PublicDashboardService) CompleteVSIncompleteProjects() (CompleteStats, error) {
	var stats CompleteStats
	if err := s.db.Model(&models.Project{}).Where("status = ?", "completed").Count(&stats.Complete).Error; err != nil {
		return stats, err
	}
	if err := s.db.Model(&models.Project{}).Where("status != ?", "completed").Count(&stats.Incomplete).Error; err != nil {
		return stats, err
	}
	return stats, nil
}

type AverageItemCostStats struct {
	AverageItemCost float64
}

func (s *PublicDashboardService) AverageItemCostPerProject() (AverageItemCostStats, error) {
	var stats AverageItemCostStats
	var totalCost float64
	var totalItems int64

	// Only consider budget items for completed projects
	var results []struct {
		ProjectID uint
		ItemCost  float64
	}
	s.db.Table("projects").
		Select("projects.id as project_id, budget_items.amount_allocated as item_cost").
		Joins("JOIN budget_items ON budget_items.project_id = projects.id").
		Where("projects.status = ?", "completed").
		Scan(&results)

	for _, r := range results {
		totalCost += r.ItemCost
		totalItems++
	}

	if totalItems > 0 {
		stats.AverageItemCost = totalCost / float64(totalItems)
	}

	return stats, nil
}

type ProjectCostVsDurationStats struct {
	AverageCostPerDay float64
}

func (s *PublicDashboardService) ProjectCostVSDuration() (ProjectCostVsDurationStats, error) {
	var stats ProjectCostVsDurationStats
	var totalCost float64
	var totalDays float64

	// Join projects and budget_items, only for completed projects
	var results []struct {
		ProjectID uint
		StartDate time.Time
		EndDate   time.Time
		TotalCost float64
	}
	s.db.Table("projects").
		Select("projects.id as project_id, projects.start_date, projects.end_date, SUM(budget_items.amount_allocated) as total_cost").
		Joins("LEFT JOIN budget_items ON budget_items.project_id = projects.id").
		Where("projects.status = ?", "completed").
		Group("projects.id, projects.start_date, projects.end_date").
		Scan(&results)

	for _, r := range results {
		days := r.EndDate.Sub(r.StartDate).Hours() / 24
		if days > 0 {
			totalCost += r.TotalCost
			totalDays += days
		}
	}

	if totalDays > 0 {
		stats.AverageCostPerDay = totalCost / totalDays
	}

	return stats, nil
}

func (s *PublicDashboardService) ProperlySpentFunds() {

}

type ProjectDurationStats struct {
	AverageEstimatedDays float64
	AverageRealDays      float64
}

func (s *PublicDashboardService) EstimatedVsRealProjectDuration() (ProjectDurationStats, error) {
	var stats ProjectDurationStats
	var totalEstimated float64
	var totalReal float64
	var count int64

	// Only consider completed projects
	var projects []struct {
		StartDate time.Time
		EndDate   time.Time
	}
	err := s.db.Model(&models.Project{}).
		Where("status = ?", "completed").
		Select("start_date, end_date").
		Scan(&projects).Error
	if err != nil {
		return stats, err
	}

	for _, p := range projects {
		days := p.EndDate.Sub(p.StartDate).Hours() / 24
		totalEstimated += days
		totalReal += days // If you have actual vs planned, adjust here
		count++
	}

	if count > 0 {
		stats.AverageEstimatedDays = totalEstimated / float64(count)
		stats.AverageRealDays = totalReal / float64(count)
	}

	return stats, nil
}

type BudgetVsDurationStats struct {
	AverageBudgetPerDay float64
}

func (s *PublicDashboardService) BudgetVsDuration() (BudgetVsDurationStats, error) {
	var stats BudgetVsDurationStats
	var totalBudget float64
	var totalDays float64

	// Join projects and budget_items, only for completed projects
	var results []struct {
		ProjectID uint
		StartDate time.Time
		EndDate   time.Time
		Budget    float64
	}
	s.db.Table("projects").
		Select("projects.id as project_id, projects.start_date, projects.end_date, SUM(budget_items.amount_allocated) as budget").
		Joins("LEFT JOIN budget_items ON budget_items.project_id = projects.id").
		Where("projects.status = ?", "completed").
		Group("projects.id, projects.start_date, projects.end_date").
		Scan(&results)

	for _, r := range results {
		days := r.EndDate.Sub(r.StartDate).Hours() / 24
		if days > 0 {
			totalBudget += r.Budget
			totalDays += days
		}
	}

	if totalDays > 0 {
		stats.AverageBudgetPerDay = totalBudget / totalDays
	}

	return stats, nil
}
