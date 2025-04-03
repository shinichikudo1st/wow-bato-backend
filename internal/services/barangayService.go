package services

import (
	"errors"
	"fmt"
	"strconv"
	database "wow-bato-backend/internal"
	"wow-bato-backend/internal/models"
)

var (
	ErrBarangayNotFound    = errors.New("barangay not found")
	ErrInvalidBarangayID   = errors.New("invalid barangay ID format")
	ErrInvalidPagination   = errors.New("invalid pagination parameters")
	ErrEmptyBarangayName   = errors.New("barangay name cannot be empty")
	ErrEmptyBarangayCity   = errors.New("barangay city cannot be empty")
	ErrEmptyBarangayRegion = errors.New("barangay region cannot be empty")
)

func validateBarangayData(barangay models.AddBarangay) error {
	if barangay.Name == "" {
		return ErrEmptyBarangayName
	}
	if barangay.City == "" {
		return ErrEmptyBarangayCity
	}
	if barangay.Region == "" {
		return ErrEmptyBarangayRegion
	}
	return nil
}


func AddNewBarangay(newBarangay models.AddBarangay) error {
	if err := validateBarangayData(newBarangay); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	db, err := database.ConnectDB()
	if err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}

	barangay := models.Barangay{
		Name:   newBarangay.Name,
		City:   newBarangay.City,
		Region: newBarangay.Region,
	}

	if result := db.Create(&barangay); result.Error != nil {
		return fmt.Errorf("failed to create barangay: %w", result.Error)
	}

	return nil
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

	if barangayUpdate.Name != "" {
		barangay.Name = barangayUpdate.Name
	}

	if barangayUpdate.City != "" {
		barangay.City = barangayUpdate.City
	}

	if barangayUpdate.Region != "" {
		barangay.Region = barangayUpdate.Region
	}

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

func OptionBarangay() ([]models.OptionBarangay, error) {
	db, err := database.ConnectDB()
	if err != nil {
		return []models.OptionBarangay{}, err
	}

	var barangay []models.OptionBarangay
	if err := db.Model(&models.Barangay{}).Select("id, name").Scan(&barangay).Error; err != nil {
		return nil, err
	}

	return barangay, nil
}

func GetSingleBarangay(barangay_ID string) (models.AllBarangayResponse, error) {
	db, err := database.ConnectDB()
	if err != nil {
		return models.AllBarangayResponse{}, fmt.Errorf("database connection failed: %w", err)
	}

	barangay_ID_int, err := strconv.Atoi(barangay_ID)
	if err != nil {
		return models.AllBarangayResponse{}, fmt.Errorf("%w: %s", ErrInvalidBarangayID, barangay_ID)
	}

	var barangay models.AllBarangayResponse
	if err := db.Model(&models.Barangay{}).
		Select("id, name, city, region").
		Where("ID = ?", barangay_ID_int).
		First(&barangay).Error; err != nil {
		return models.AllBarangayResponse{}, fmt.Errorf("%w: ID %d", ErrBarangayNotFound, barangay_ID_int)
	}

	return barangay, nil
}

func AllBarangaysPublic() ([]models.PublicBarangayDisplay, error) {
	db, err := database.ConnectDB()
	if err != nil {
		return []models.PublicBarangayDisplay{}, err
	}

	var barangays []models.PublicBarangayDisplay
	if err := db.Find(&barangays).Error; err != nil {
		return []models.PublicBarangayDisplay{}, err
	}

	return barangays, nil
}
