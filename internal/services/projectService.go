package services

import (
	"errors"
	"fmt"
	"strconv"
	"time"
	"wow-bato-backend/internal/models"

	"gorm.io/gorm"
)

var (
	GO_DATE_FORMAT = "2006-01-02"
	ErrParseStartDate = errors.New("something went wrong while parsing start date")  
	ErrParseEndDate = errors.New("something went wrong while parsing end date") 
	ErrProjectUpdate = errors.New("project to update not found")	
)


type ProjectService struct {
	db *gorm.DB
}

func NewProjectService(db *gorm.DB) *ProjectService {
	return &ProjectService{db: db}
}

func (s *ProjectService) AddNewProject(barangay_ID uint, categoryID string, newProject models.NewProject) error {

	categoryID_int, err := strconv.Atoi(categoryID)
	if err != nil {
		return err
	}

	layout := GO_DATE_FORMAT

	startDate, err := time.Parse(layout, newProject.StartDate)
	if err != nil {
		return ErrParseStartDate
	}

	endDate, err := time.Parse(layout, newProject.EndDate)
	if err != nil {
		return ErrParseEndDate
	}

	project := models.Project{
		Barangay_ID: barangay_ID,
		CategoryID:  uint(categoryID_int),
		Name:        newProject.Name,
		Description: newProject.Description,
		StartDate:   startDate,
		EndDate:     endDate,
		Status:      newProject.Status,
	}

	result := s.db.Create(&project)

	return result.Error
}

func (s *ProjectService) DeleteProject(barangay_ID uint, projectID string) error {

	projectID_int, err := strconv.Atoi(projectID)
	if err != nil {
		return err
	}

	result := s.db.Delete(&models.Project{}, projectID_int)

	return result.Error
}

func (s *ProjectService) UpdateProject(barangay_ID uint, projectID string, updateProject models.UpdateProject) error {

	projectID_int, err := strconv.Atoi(projectID)
	if err != nil {
		return err
	}

	var project models.Project
	if err := s.db.Where("Barangay_ID = ? AND id = ?", barangay_ID, projectID_int).Error; err != nil {
		return fmt.Errorf("update failed: %w", ErrProjectUpdate)
	}

	project.Name = updateProject.Name
	project.Description = updateProject.Description

	result := s.db.Save(&project)

	return result.Error
}

func (s *ProjectService) GetAllProjects(barangay_ID uint, categoryID string, limit string, page string) ([]models.ProjectList, error) {

	limit_int, err := strconv.Atoi(limit)
	if err != nil {
		return []models.ProjectList{}, err
	}

	page_int, err := strconv.Atoi(page)
	if err != nil {
		return []models.ProjectList{}, err
	}

	offset := (page_int - 1) * limit_int

	categoryID_int, err := strconv.Atoi(categoryID)
	if err != nil {
		return []models.ProjectList{}, err
	}

	var projects []models.ProjectList
	if err := s.db.Model(&models.Project{}).
		Where("barangay_id = ? AND category_id = ?", barangay_ID, categoryID_int).
		Select("id, name, status, start_date, end_date").
		Limit(limit_int).
		Offset(offset).
		Scan(&projects).Error; err != nil {
		return []models.ProjectList{}, fmt.Errorf("failed to retrieve all projects: %w", err)
	}

	return projects, nil
}

func (s *ProjectService) UpdateProjectStatus(projectID string, barangay_ID uint, newStatus models.NewProjectStatus) error {

	projectID_int, err := strconv.Atoi(projectID)
	if err != nil {
		return err
	}

	var project models.Project
	if err := s.db.Where("Barangay_ID = ? AND id = ?", barangay_ID, projectID_int).First(&project).Error; err != nil {
		return fmt.Errorf("project not found: %w", err)
	}

	if newStatus.Status == "ongoing" {

		project.Status = newStatus.Status
		project.StartDate = newStatus.FlexDate

		result := s.db.Save(&project)
		return result.Error

	} else {
		project.Status = newStatus.Status
		project.EndDate = newStatus.FlexDate

		result := s.db.Save(&project)
		return result.Error

	}
}

func (s *ProjectService) GetProjectSingle(projectID string)(models.ProjectList, error){

	projectID_int, err := strconv.Atoi(projectID)
	if err != nil {
		return models.ProjectList{}, err
	}

	var project models.ProjectList
	if err := s.db.Model(&models.Project{}).Where("id = ?", projectID_int).First(&project).Error; err != nil {
		return models.ProjectList{}, fmt.Errorf("project not found: %w", err)
	}

	return project, nil
}