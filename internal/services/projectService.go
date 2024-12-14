package services

import (
	"strconv"
	"time"
	database "wow-bato-backend/internal"
	"wow-bato-backend/internal/models"
)

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
		CategoryID: uint(categoryID_int),
		Name: newProject.Name,
		Description: newProject.Description,
		StartDate: startDate,
        EndDate: endDate,
		Status: newProject.Status,
	}

	result := db.Create(&project)

	return result.Error
}

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

func GetAllProjects(barangay_ID uint, categoryID string) ([]models.ProjectList, error) {
	db, err := database.ConnectDB()
	if err != nil {
		return nil, err
	}

    categoryID_int, err := strconv.Atoi(categoryID)
    if err != nil {
        return []models.ProjectList{}, err
    }

	var projects []models.ProjectList
    if err := db.Model(&models.Project{}).
        Where("barangay_id = ? AND category_id = ?", barangay_ID, categoryID_int).
        Select("id, name, status, start_date, end_date").
        Scan(&projects).Error; err != nil {
            return []models.ProjectList{}, err
        }

	return projects, nil
}

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
