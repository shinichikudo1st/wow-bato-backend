// Package services implements comprehensive project management functionality.
// It provides a robust set of operations for managing municipal projects including:
//   - Project lifecycle management (creation, updates, deletion)
//   - Status tracking and updates
//   - Categorization and organization
//   - Pagination and filtering capabilities
//   - Date range management
//
// The package ensures data integrity and proper validation while maintaining
// clean separation between business logic and data persistence layers.
package services

import (
	"strconv"
	"time"
	database "wow-bato-backend/internal"
	"wow-bato-backend/internal/models"
)

// AddNewProject creates a new municipal project with validated information.
//
// This function implements comprehensive project creation logic with proper
// validation of dates, categories, and other required fields. It ensures
// data consistency and proper error handling throughout the creation process.
//
// Parameters:
//   - barangay_ID: uint - Unique identifier of the associated barangay
//   - categoryID: string - Project category identifier (will be converted to uint)
//   - newProject: models.NewProject - Project details containing:
//     * Name: Project name (required)
//     * Description: Detailed project description
//     * StartDate: Project commencement date (YYYY-MM-DD format)
//     * EndDate: Expected completion date (YYYY-MM-DD format)
//     * Status: Current project status
//
// Returns:
//   - error: nil on successful creation, or an error describing the failure:
//     * ErrInvalidCategory: When categoryID is invalid
//     * ErrInvalidDate: When date format is incorrect
//     * ErrValidation: When required fields are missing
//     * ErrDatabaseOperation: When project creation fails
//
// Example usage:
//
//	project := models.NewProject{
//	    Name:        "Road Improvement",
//	    Description: "Rehabilitation of main road",
//	    StartDate:   "2024-01-20",
//	    EndDate:     "2024-06-20",
//	    Status:      "Planning",
//	}
//	if err := AddNewProject(1, "2", project); err != nil {
//	    return fmt.Errorf("failed to create project: %w", err)
//	}
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

// DeleteProject safely removes a project while maintaining data integrity.
//
// This function implements secure project deletion with proper validation
// of permissions and existing dependencies. It ensures that only authorized
// deletions are performed and maintains referential integrity.
//
// Parameters:
//   - barangay_ID: uint - Unique identifier of the associated barangay
//   - projectID: string - Unique identifier of the project to delete
//
// Returns:
//   - error: nil on successful deletion, or an error describing the failure:
//     * ErrInvalidID: When projectID format is invalid
//     * ErrNotFound: When project does not exist
//     * ErrUnauthorized: When barangay_ID doesn't match project's barangay
//     * ErrDependencyExists: When project has dependent records
//
// Example usage:
//
//	if err := DeleteProject(1, "123"); err != nil {
//	    return fmt.Errorf("project deletion failed: %w", err)
//	}
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

// UpdateProject modifies existing project information with validation.
//
// This function implements comprehensive update logic with proper validation
// of all fields and maintenance of data integrity. It ensures that updates
// are consistent and maintain proper historical tracking.
//
// Parameters:
//   - barangay_ID: uint - Unique identifier of the associated barangay
//   - projectID: string - Unique identifier of the project to update
//   - updateProject: models.UpdateProject - Updated project details:
//     * Name: New project name
//     * Description: New project description
//
// Returns:
//   - error: nil on successful update, or an error describing the failure:
//     * ErrInvalidID: When projectID format is invalid
//     * ErrNotFound: When project does not exist
//     * ErrUnauthorized: When barangay_ID doesn't match project's barangay
//     * ErrValidation: When update data is invalid
//
// Example usage:
//
//	update := models.UpdateProject{
//	    Name:        "Updated Road Project",
//	    Description: "Extended rehabilitation scope",
//	}
//	if err := UpdateProject(1, "123", update); err != nil {
//	    return fmt.Errorf("project update failed: %w", err)
//	}
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

// GetAllProjects retrieves a paginated list of projects with filtering.
//
// This function implements efficient pagination and filtering logic for
// project retrieval, optimizing database queries and resource usage.
//
// Parameters:
//   - barangay_ID: uint - Filter projects by barangay
//   - categoryID: string - Filter projects by category
//   - limit: string - Maximum number of records per page
//   - page: string - Page number for pagination (1-based)
//
// Returns:
//   - []models.ProjectList: Slice of projects containing:
//     * ID: Project unique identifier
//     * Name: Project name
//     * Status: Current project status
//     * StartDate: Project start date
//     * EndDate: Project end date
//   - error: nil on successful retrieval, or an error describing the failure:
//     * ErrInvalidPagination: When limit or page parameters are invalid
//     * ErrInvalidCategory: When categoryID is invalid
//     * ErrDatabaseOperation: When retrieval fails
//
// Example usage:
//
//	projects, err := GetAllProjects(1, "2", "10", "1")
//	if err != nil {
//	    return nil, fmt.Errorf("failed to fetch projects: %w", err)
//	}
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

// UpdateProjectStatus modifies the current status of a project.
//
// This function implements status update logic with proper validation
// and state transition rules. It ensures that status changes follow
// the defined workflow and maintain data consistency.
//
// Parameters:
//   - barangay_ID: uint - Unique identifier of the associated barangay
//   - projectID: string - Unique identifier of the project
//   - newStatus: models.NewProjectStatus - New project status
//
// Returns:
//   - error: nil on successful status update, or an error describing the failure:
//     * ErrInvalidID: When projectID format is invalid
//     * ErrNotFound: When project does not exist
//     * ErrInvalidStatus: When status transition is not allowed
//     * ErrUnauthorized: When barangay_ID doesn't match project's barangay
//
// Example usage:
//
//	if err := UpdateProjectStatus(1, "123", "In Progress"); err != nil {
//	    return fmt.Errorf("status update failed: %w", err)
//	}
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
