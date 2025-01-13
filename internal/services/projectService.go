// Package services provides project-related business logic and operations for the application.
// It handles project management while maintaining separation of concerns from the database and presentation layers.
// It does the following:
// - AddNewProject: Adds a new project to the database.
// - DeleteProject: Deletes a project from the database.
// - UpdateProject: Updates a project in the database.
// - GetAllProjects: Retrieves a list of projects from the database.
// - UpdateProjectStatus: Updates the status of a project in the database.
package services

import (
	"strconv"
	"time"
	database "wow-bato-backend/internal"
	"wow-bato-backend/internal/models"
)

// AddNewProject adds a new project to the database
//
// This function does the following operations:
// 1. Establishes a database connection
// 2. Converts the categoryID from string to int
// 3. Converts the startDate and endDate from string to time
// 4. Creates a new project record in the database
//
// Parameters:
//   - barangay_ID: uint - The barangay ID of the project
//   - categoryID: string - The category ID of the project
//   - newProject: models.NewProject - The new project data
//
// Returns:
//   - error: Returns nil if successful, otherwise returns an error:
//   - Database connection errors
//   - Database creation errors
func AddNewProject(barangay_ID uint, categoryID string, newProject models.NewProject) error {
	db, err := database.ConnectDB()
	if err != nil {
		return err
	}

	categoryID_int, err := strconv.Atoi(categoryID)
	if err != nil {
		return err
	}

	layout := "2006-01-02"

	startDate, err := time.Parse(layout, newProject.StartDate)
	if err != nil {
		return err
	}

	endDate, err := time.Parse(layout, newProject.EndDate)
	if err != nil {
		return err
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

	result := db.Create(&project)

	return result.Error
}

// DeleteProject deletes a project from the database
//
// This function does the following operations:
// 1. Establishes a database connection
// 2. Converts the projectID from string to int
// 3. Deletes the project record from the database
//
// Parameters:
//   - barangay_ID: uint - The barangay ID of the project
//   - projectID: string - The ID of the project to delete
//
// Returns:
//   - error: Returns nil if successful, otherwise returns an error:
//   - Database connection errors
//   - Database deletion errors
func DeleteProject(barangay_ID uint, projectID string) error {
	db, err := database.ConnectDB()
	if err != nil {
		return err
	}

	projectID_int, err := strconv.Atoi(projectID)
	if err != nil {
		return err
	}

	result := db.Delete(&models.Project{}, projectID_int)

	return result.Error
}

// UpdateProject updates a project in the database
//
// This function does the following operations:
// 1. Establishes a database connection
// 2. Converts the projectID from string to int
// 3. Updates the project record in the database
//
// Parameters:
//   - barangay_ID: uint - The barangay ID of the project
//   - projectID: string - The ID of the project to update
//   - updateProject: models.UpdateProject - The updated project data
//
// Returns:
//   - error: Returns nil if successful, otherwise returns an error:
//   - Database connection errors
//   - Database update errors
func UpdateProject(barangay_ID uint, projectID string, updateProject models.UpdateProject) error {
	db, err := database.ConnectDB()
	if err != nil {
		return err
	}

	projectID_int, err := strconv.Atoi(projectID)
	if err != nil {
		return err
	}

	var project models.Project
	if err := db.Where("Barangay_ID = ? AND id = ?", barangay_ID, projectID_int).Error; err != nil {
		return err
	}

	project.Name = updateProject.Name
	project.Description = updateProject.Description

	result := db.Save(&project)

	return result.Error
}

// GetAllProjects retrieves a list of projects from the database
//
// This function does the following operations:
// 1. Establishes a database connection
// 2. Converts the limit and page from string to int
// 3. Retrieves a list of projects from the database
//
// Parameters:
//   - barangay_ID: uint - The barangay ID of the project
//   - categoryID: string - The category ID of the project
//   - limit: string - The number of projects to retrieve per page
//   - page: string - The page number to retrieve
//
// Returns:
//   - []models.ProjectList: A slice of project responses
//   - error: Returns nil if successful, otherwise returns an error:
//   - Database connection errors
func GetAllProjects(barangay_ID uint, categoryID string, limit string, page string) ([]models.ProjectList, error) {
	db, err := database.ConnectDB()
	if err != nil {
		return nil, err
	}

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
	if err := db.Model(&models.Project{}).
		Where("barangay_id = ? AND category_id = ?", barangay_ID, categoryID_int).
		Select("id, name, status, start_date, end_date").
		Limit(limit_int).
		Offset(offset).
		Scan(&projects).Error; err != nil {
		return []models.ProjectList{}, err
	}

	return projects, nil
}

// UpdateProjectStatus updates the status of a project in the database
//
// This function does the following operations:
// 1. Establishes a database connection
// 2. Converts the projectID from string to int
// 3. Updates the project status in the database
//
// Parameters:
//   - projectID: string - The ID of the project to update
//   - barangay_ID: uint - The barangay ID of the project
//   - newStatus: models.NewProjectStatus - The new status of the project
//
// Returns:
//   - error: Returns nil if successful, otherwise returns an error:
//   - Database connection errors
func UpdateProjectStatus(projectID string, barangay_ID uint, newStatus models.NewProjectStatus) error {
	db, err := database.ConnectDB()
	if err != nil {
		return err
	}

	projectID_int, err := strconv.Atoi(projectID)
	if err != nil {
		return err
	}

	var project models.Project
	if err := db.Where("Barangay_ID = ? AND id = ?", barangay_ID, projectID_int).First(&project).Error; err != nil {
		return err
	}

	if newStatus.Status == "ongoing" {

		project.Status = newStatus.Status
		project.StartDate = newStatus.FlexDate

		result := db.Save(&project)
		return result.Error

	} else {
		project.Status = newStatus.Status
		project.EndDate = newStatus.FlexDate

		result := db.Save(&project)
		return result.Error

	}
}
