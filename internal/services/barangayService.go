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

func UpdateBarangay(barangay_ID string, barangayUpdate models.UpdateBarangay) error {
	db, err := database.ConnectDB()
	if err != nil {
		return err
	}

	barangay_ID_int, err := strconv.Atoi(barangay_ID)
	if err != nil {
		return err
	}

	var barangay models.Barangay

	if err := db.Where("id = ?", barangay_ID_int).First(&barangay).Error; err != nil {
		return err
	}

	barangay.Name = barangayUpdate.Name
	barangay.City = barangayUpdate.City
	barangay.Region = barangayUpdate.Region

	result := db.Save(&barangay)

	return result.Error
}

func GetAllBarangay(limit string, page string) ([]models.AllBarangayResponse, error) {
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
	
	offset := (pageInt - 1) * limitInt

	var barangay []models.AllBarangayResponse
	if err := db.Model(&models.Barangay{}).Select("id, name, city, region").Limit(limitInt).Offset(offset).Find(&barangay).Error; err != nil {
		return nil, err
	}

	return barangay, nil
}

func GetSingleBarangay(id string)(models.AllBarangayResponse, error){
	db, err := database.ConnectDB()
	if err != nil {
		return models.AllBarangayResponse{}, err
	}

	barangay_ID, err := strconv.Atoi(id)
	if err != nil {
		return models.AllBarangayResponse{}, err
	}

	var barangay models.AllBarangayResponse
    if err := db.Model(&models.Barangay{}).Select("id, name, city, region").Where("ID = ?", barangay_ID).First(&barangay).Error; err != nil {
		return models.AllBarangayResponse{}, err
	}

	return barangay, nil
}
