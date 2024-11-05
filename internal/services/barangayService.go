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
		Name:   newBarangay.Name,
		City:   newBarangay.City,
		Region: newBarangay.Region,
	}

	result := db.Create(&barangay)

	return result.Error
}

func DeleteBarangay(barangayToDelete models.DeleteBarangay) error {
	db, err := database.ConnectDB()

	if err != nil {
		return err
	}

	var barangay models.Barangay
	result := db.Where("Barangay_ID = ?", barangayToDelete.Barangay_ID).Delete(&barangay)

	return result.Error
}

func UpdateBarangay(barangayToUpdate models.UpdateBarangay) error {
	db, err := database.ConnectDB()
	if err != nil {
		return err
	}

	var barangay models.Barangay

	if err := db.Where("Barangay_ID = ?", barangayToUpdate.Barangay_ID).First(&barangay).Error; err != nil {
		return err
	}

	barangay.Name = barangayToUpdate.Name
	barangay.City = barangayToUpdate.City
	barangay.Region = barangayToUpdate.Region

	result := db.Save(&barangay)

	return result.Error
}

func GetAllBarangay() ([]models.Barangay, error) {
	db, err := database.ConnectDB()
	if err != nil {
		return nil, err
	}

	var barangay []models.Barangay
	result := db.Find(&barangay)

	return barangay, result.Error
}
