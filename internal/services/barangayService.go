package services

import (
	database "wow-bato-backend/internal"
	"wow-bato-backend/internal/models"
)

func AddNewBarangay(newBarangay models.AddBarangay) error {
	db, err := database.ConnectDB()
	if err != nil {
		return err
	}

	barangay := models.Barangay{
		Name: newBarangay.Name,
		City: newBarangay.City,
		Region: newBarangay.Region,
	}

	result := db.Create(&barangay)

	return result.Error
}