package services

import (
	"strconv"
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

	project := models.Project{
		Barangay_ID: barangay_ID,
		CategoryID: uint(categoryID_int),
		Name: newProject.Name,
		Description: newProject.Description,
		StartDate: newProject.StartDate,
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

func GetAllProjects(barangay_ID uint) ([]models.Project, error) {
	db, err := database.ConnectDB()
	if err != nil {
		return nil, err
	}

	var projects []models.Project
	result := db.Where("barangay_id = ?", barangay_ID).Find(&projects)

	return projects, result.Error
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
        project.EndDate = &newStatus.FlexDate

        result := db.Save(&project)
        return result.Error

    }
}
