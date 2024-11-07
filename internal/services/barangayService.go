package services

import (
	"strconv"
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

func DeleteBarangay(barangay_ID string) error {
	db, err := database.ConnectDB()
	if err != nil {
		return err
	}

	barangay_ID_int, err := strconv.Atoi(barangay_ID)
	if err != nil {
		return err
	}

	var barangay models.Barangay
	result := db.Where("id = ?", barangay_ID_int).Delete(&barangay)

	return result.Error
}

func UpdateBarangay(barangayToUpdate models.UpdateBarangay) error {
	db, err := database.ConnectDB()
	if err != nil {
		return err
	}

	var barangay models.Barangay

	if err := db.Where("id = ?", barangayToUpdate.Barangay_ID).First(&barangay).Error; err != nil {
		return err
	}

	barangay.Name = barangayToUpdate.Name
	barangay.City = barangayToUpdate.City
	barangay.Region = barangayToUpdate.Region

	result := db.Save(&barangay)

	return result.Error
}

func GetAllBarangay(limit string, page string) ([]models.Barangay, error) {
	db, err := database.ConnectDB()
	if err != nil {
		return nil, err
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return nil, err
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return nil, err
	}

	var barangay []models.Barangay
	result := db.Limit(limitInt).Offset(pageInt).Find(&barangay)

	return barangay, result.Error
}

func GetSingleBarangay(id string)(models.Barangay, error){
	db, err := database.ConnectDB()
	if err != nil {
		return models.Barangay{}, err
	}

	barangay_ID, err := strconv.Atoi(id)
	if err != nil {
		return models.Barangay{}, err
	}

	var barangay models.Barangay
	if err := db.Where("id = ?", barangay_ID).First(&barangay); err != nil {
		return models.Barangay{}, err.Error
	}

	return barangay, nil
}
