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
