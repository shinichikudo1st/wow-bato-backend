package services

import (
	"errors"
	"fmt"
	"strconv"
	"wow-bato-backend/internal/models"

	"gorm.io/gorm"
)

var (
	ErrBarangayNotFound    = errors.New("barangay not found")
	ErrInvalidBarangayID   = errors.New("invalid barangay ID format")
	ErrInvalidPagination   = errors.New("invalid pagination parameters")
	ErrEmptyBarangayName   = errors.New("barangay name cannot be empty")
	ErrEmptyBarangayCity   = errors.New("barangay city cannot be empty")
	ErrEmptyBarangayRegion = errors.New("barangay region cannot be empty")
)

type BarangayService struct {
	db *gorm.DB
}

func NewBarangayService(db *gorm.DB) *BarangayService {
	return &BarangayService{db: db}
}

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

func (s *BarangayService) AddNewBarangay(newBarangay models.AddBarangay) error {
	if err := validateBarangayData(newBarangay); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	barangay := models.Barangay{
		Name:   newBarangay.Name,
		City:   newBarangay.City,
		Region: newBarangay.Region,
	}

	if err := s.db.Create(&barangay).Error; err != nil {
		return fmt.Errorf("failed to create barangay: %w", err)
	}

	return nil
}

func (s *BarangayService) DeleteBarangay(barangay_ID string) error {

	barangay_ID_int, err := strconv.Atoi(barangay_ID)
	if err != nil {
		return err
	}

	var barangay models.Barangay
	if err := s.db.Where("id = ?", barangay_ID_int).Delete(&barangay).Error; err != nil {
		return fmt.Errorf("failed to delete barangay: %w", err)
	}

	return nil
}

func (s *BarangayService) UpdateBarangay(barangay_ID string, barangayUpdate models.UpdateBarangay) error {

	barangay_ID_int, err := strconv.Atoi(barangay_ID)
	if err != nil {
		return err
	}

	var barangay models.Barangay

	if err := s.db.Where("id = ?", barangay_ID_int).First(&barangay).Error; err != nil {
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

	if err := s.db.Save(&barangay).Error; err != nil {
		return fmt.Errorf("failed to save new barangay info changes: %w", err)
	}

	return nil
}

func (s *BarangayService) GetAllBarangay(limit string, page string) ([]models.AllBarangayResponse, error) {

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
	if err := s.db.Model(&models.Barangay{}).
		Select("id, name, city, region").
		Limit(limitInt).Offset(offset).
		Find(&barangay).Error; 
		err != nil { return nil, err }

	return barangay, nil
}

func (s *BarangayService) OptionBarangay() ([]models.OptionBarangay, error) {

	var barangay []models.OptionBarangay
	if err := s.db.Model(&models.Barangay{}).
		Select("id, name").Scan(&barangay).Error; 
		err != nil {
			return []models.OptionBarangay{}, err
		}

	return barangay, nil
}

func (s *BarangayService) GetSingleBarangay(barangay_ID string) (models.AllBarangayResponse, error) {

	barangay_ID_int, err := strconv.Atoi(barangay_ID)
	if err != nil {
		return models.AllBarangayResponse{}, fmt.Errorf("%w: %s", ErrInvalidBarangayID, barangay_ID)
	}

	var barangay models.AllBarangayResponse
	if err := s.db.Model(&models.Barangay{}).
		Select("id, name, city, region").
		Where("ID = ?", barangay_ID_int).
		First(&barangay).Error; err != nil {
		return models.AllBarangayResponse{}, fmt.Errorf("%w: ID %d", ErrBarangayNotFound, barangay_ID_int)
	}

	return barangay, nil
}

func (s *BarangayService) AllBarangaysPublic() ([]models.PublicBarangayDisplay, error) {

	var barangays []models.PublicBarangayDisplay
	if err := s.db.Find(&barangays).Error; err != nil {
		return nil, err
	}

	return barangays, nil
}
