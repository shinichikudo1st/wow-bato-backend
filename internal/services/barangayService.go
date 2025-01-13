// The barangayService package provides functions for managing barangays in the application.
// It includes functions for adding new barangays, deleting barangays, updating barangays,
// retrieving all barangays, and retrieving a single barangay.
package services

import (
	"strconv"
	database "wow-bato-backend/internal"
	"wow-bato-backend/internal/models"
)

// AddNewBarangay adds a new barangay to the database.
//
// This function performs the following operations:
// 1. Establishes a database connection
// 2. Creates a new barangay record in the database
//
// Parameters:
//   - newBarangay: models.AddBarangay -
//     Contains barangay data including:
//   - Name: Barangay name
//   - City: Barangay city
//   - Region: Barangay region
//
// Returns:
//   - error: Returns nil if successful, otherwise returns an error:
//   - Database connection errors
//   - Database creation errors
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

// DeleteBarangay deletes a barangay from the database.
//
// This function performs the following operations:
//  1. Establishes a database connection
//  2. Converts the barangay ID from string to int
//     because we want the handler to be clean and secure
//  3. Deletes the barangay record from the database that matches the barangay_ID
//
// Parameters:
//   - barangay_ID: string - The ID of the barangay to be deleted
//
// Returns:
//   - error: Returns nil if successful, otherwise returns an error:
//   - Database connection errors
//   - Database deletion errors
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

// UpdateBarangay updates a barangay in the database.
//
// This function performs the following operations:
//  1. Establishes a database connection
//  2. Converts the barangay ID from string to int
//     because we want the handler to be clean and secure
//  3. Updates the barangay record in the database that matches the barangay_ID
//
// Parameters:
//   - barangay_ID: string - The ID of the barangay to be updated
//   - barangayUpdate: models.UpdateBarangay -
//     Contains barangay data including:
//   - Name: Barangay name
//   - City: Barangay city
//   - Region: Barangay region
//
// Returns:
//   - error: Returns nil if successful, otherwise returns an error:
//   - Database connection errors
//   - Database update errors
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

// GetAllBarangay retrieves all barangays from the database.
//
// This function performs the following operations:
//  1. Establishes a database connection
//  2. Converts the limit and page from string to int
//  3. Retrieves all barangays from the database
//  4. Returns the barangays and any errors
//
// Parameters:
//   - limit: string - The number of barangays to retrieve per page
//   - page: string - The page number to retrieve
//
// Returns:
//   - []models.AllBarangayResponse: A slice of barangay responses
//   - error: Returns nil if successful, otherwise returns an error:
//   - Database connection errors
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

// OptionBarangay retrieves all barangays from the database.
// This is only used for the barangay dropdown selection during user creation.
//
// This function performs the following operations:
//  1. Establishes a database connection
//  2. Retrieves all barangays from the database selecting id and name
//  3. Returns the barangays and any errors
//
// Returns:
//   - []models.OptionBarangay: A slice of barangay options
//   - error: Returns nil if successful, otherwise returns an error:
//   - Database connection errors
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

// GetSingleBarangay retrieves a single barangay from the database.
//
// This function performs the following operations:
//  1. Establishes a database connection
//  2. Converts the barangay ID from string to int
//     because we want the handler to be clean and secure
//  3. Retrieves the barangay record from the database that matches the barangay_ID
//  4. Returns the barangay and any errors
//
// Parameters:
//   - barangay_ID: string - The ID of the barangay to be retrieved
//
// Returns:
//   - models.AllBarangayResponse: A barangay response
//   - error: Returns nil if successful, otherwise returns an error:
//   - Database connection errors
func GetSingleBarangay(id string) (models.AllBarangayResponse, error) {
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
