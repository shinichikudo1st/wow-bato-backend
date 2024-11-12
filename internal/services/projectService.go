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
