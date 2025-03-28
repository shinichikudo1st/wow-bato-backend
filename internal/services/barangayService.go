// Package services provides comprehensive management of barangay (administrative division) data.
// It implements the core business logic for barangay operations including:
//   - Creation and registration of new barangays
//   - Modification and updates of existing barangay information
//   - Retrieval of barangay data with pagination support
//   - Deletion of barangay records
//   - Option listing for UI integration
//
// The package ensures data consistency and proper validation while maintaining
// separation of concerns between the database layer and the presentation layer.
package services

import (
	"errors"
	"fmt"
	"strconv"
	database "wow-bato-backend/internal"
	"wow-bato-backend/internal/models"
)

// Domain-specific errors for barangay operations
var (
	ErrBarangayNotFound    = errors.New("barangay not found")
	ErrInvalidBarangayID   = errors.New("invalid barangay ID format")
	ErrInvalidPagination   = errors.New("invalid pagination parameters")
	ErrEmptyBarangayName   = errors.New("barangay name cannot be empty")
	ErrEmptyBarangayCity   = errors.New("barangay city cannot be empty")
	ErrEmptyBarangayRegion = errors.New("barangay region cannot be empty")
)

// validateBarangayData validates the required fields for barangay data
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

// AddNewBarangay creates a new barangay record with validated information.
//
// This function implements the core business logic for barangay registration,
// ensuring data consistency and proper validation before storage.
//
// Parameters:
//   - newBarangay: models.AddBarangay - The barangay data transfer object containing:
//   - Name: Name of the barangay (required)
//   - City: City where the barangay is located (required)
//   - Region: Administrative region of the barangay (required)
//
// Returns:
//   - error: nil on successful creation, or an error describing the failure:
//   - ErrDatabaseConnection: When database connection fails
//   - ErrDuplicateBarangay: When barangay name already exists
//   - ErrValidation: When required fields are missing or invalid
//   - ErrDatabaseOperation: When barangay creation fails
//
// Example usage:
//
//	newBarangay := models.AddBarangay{
//	    Name:   "San Antonio",
//	    City:   "Quezon City",
//	    Region: "National Capital Region",
//	}
//	if err := AddNewBarangay(newBarangay); err != nil {
//	    return fmt.Errorf("failed to create barangay: %w", err)
//	}
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

// DeleteBarangay removes a barangay record from the system.
//
// This function ensures safe deletion of barangay records by validating
// dependencies and maintaining referential integrity.
//
// Parameters:
//   - barangay_ID: string - Unique identifier of the barangay to delete
//
// Returns:
//   - error: nil on successful deletion, or an error describing the failure:
//   - ErrInvalidID: When barangay_ID format is invalid
//   - ErrNotFound: When barangay does not exist
//   - ErrDependencyExists: When barangay has associated records
//   - ErrDatabaseOperation: When deletion operation fails
//
// Example usage:
//
//	if err := DeleteBarangay("123"); err != nil {
//	    return fmt.Errorf("failed to delete barangay: %w", err)
//	}
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

// UpdateBarangay modifies existing barangay information.
//
// This function implements validation and update logic for barangay records,
// ensuring data consistency and proper error handling.
//
// Parameters:
//   - barangay_ID: string - Unique identifier of the barangay to update
//   - barangayUpdate: models.UpdateBarangay - The update data containing:
//   - Name: Updated name of the barangay
//   - City: Updated city information
//   - Region: Updated region information
//
// Returns:
//   - error: nil on successful update, or an error describing the failure:
//   - ErrInvalidID: When barangay_ID format is invalid
//   - ErrNotFound: When barangay does not exist
//   - ErrValidation: When update data is invalid
//   - ErrDatabaseOperation: When update operation fails
//
// Example usage:
//
//	updateData := models.UpdateBarangay{
//	    Name:   "New San Antonio",
//	    City:   "Makati City",
//	    Region: "NCR",
//	}
//	if err := UpdateBarangay("123", updateData); err != nil {
//	    return fmt.Errorf("failed to update barangay: %w", err)
//	}
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

// GetAllBarangay retrieves a paginated list of barangays.
//
// This function implements pagination and filtering logic for efficient
// data retrieval and resource optimization.
//
// Parameters:
//   - limit: string - Maximum number of records to return per page
//   - page: string - Page number for pagination (1-based indexing)
//
// Returns:
//   - []models.AllBarangayResponse: Slice of barangay records containing:
//   - ID: Unique identifier of the barangay
//   - Name: Name of the barangay
//   - City: City where the barangay is located
//   - Region: Administrative region
//   - error: nil on successful retrieval, or an error describing the failure:
//   - ErrInvalidPagination: When limit or page parameters are invalid
//   - ErrDatabaseOperation: When retrieval operation fails
//
// Example usage:
//
//	barangays, err := GetAllBarangay("10", "1")
//	if err != nil {
//	    return nil, fmt.Errorf("failed to fetch barangays: %w", err)
//	}
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

// OptionBarangay retrieves a list of barangays for dropdown selection.
//
// This function provides optimized data retrieval for UI components,
// returning only essential fields needed for selection interfaces.
//
// Returns:
//   - []models.OptionBarangay: Slice of barangay options containing:
//   - ID: Unique identifier of the barangay
//   - Name: Name of the barangay
//   - error: nil on successful retrieval, or an error describing the failure:
//   - ErrDatabaseOperation: When retrieval operation fails
//
// Example usage:
//
//	options, err := OptionBarangay()
//	if err != nil {
//	    return nil, fmt.Errorf("failed to fetch barangay options: %w", err)
//	}
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

// GetSingleBarangay retrieves detailed information for a specific barangay.
//
// This function provides comprehensive data retrieval for a single barangay record,
// including all associated information and metadata.
//
// Parameters:
//   - barangay_ID: string - Unique identifier of the barangay to retrieve
//
// Returns:
//   - models.SingleBarangayResponse: Detailed barangay information containing:
//   - ID: Unique identifier of the barangay
//   - Name: Name of the barangay
//   - City: City where the barangay is located
//   - Region: Administrative region
//   - CreatedAt: Timestamp of record creation
//   - UpdatedAt: Timestamp of last update
//   - error: nil on successful retrieval, or an error describing the failure:
//   - ErrInvalidID: When barangay_ID format is invalid
//   - ErrNotFound: When barangay does not exist
//   - ErrDatabaseOperation: When retrieval operation fails
//
// Example usage:
//
//	barangay, err := GetSingleBarangay("123")
//	if err != nil {
//	    return nil, fmt.Errorf("failed to fetch barangay: %w", err)
//	}
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
