package services

import (
	"errors"
	"fmt"
	"strconv"
	"wow-bato-backend/internal/models"

	"gorm.io/gorm"
)

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrInvalidCredentials    = errors.New("invalid email or password")
	ErrEmptyEmail            = errors.New("email cannot be empty")
	ErrEmptyPassword         = errors.New("password cannot be empty")
	ErrEmptyFirstName        = errors.New("first name cannot be empty")
	ErrEmptyLastName         = errors.New("last name cannot be empty")
	ErrInvalidRole           = errors.New("invalid user role")
	ErrEmptyContact          = errors.New("contact information cannot be empty")
	ErrPasswordHashingFailed = errors.New("password hashing failed")
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func validateUserRegistration(user models.RegisterUser) error {
	if user.Email == "" {
		return ErrEmptyEmail
	}
	if user.Password == "" {
		return ErrEmptyPassword
	}
	if user.FirstName == "" {
		return ErrEmptyFirstName
	}
	if user.LastName == "" {
		return ErrEmptyLastName
	}
	if user.Contact == "" {
		return ErrEmptyContact
	}

	return nil
}

func validateLoginUser(loginUser models.LoginUser) (models.UserStruct, error){
	if loginUser.Email == "" {
		return models.UserStruct{}, ErrEmptyEmail
	}
	if loginUser.Password == "" {
		return models.UserStruct{}, ErrEmptyPassword
	}

	return models.UserStruct{}, nil
}

func (s *UserService) RegisterUser(registerUser models.RegisterUser) error {
	if err := validateUserRegistration(registerUser); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	hash, err := HashPassword(registerUser.Password)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrPasswordHashingFailed, err)
	}

	barangay_ID, err := strconv.Atoi(registerUser.Barangay_ID)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrorInvalidBarangayID, registerUser.Barangay_ID)
	}

	barangay_ID_uint := uint(barangay_ID)

	user := models.User{
		Email:       registerUser.Email,
		Password:    hash,
		FirstName:   registerUser.FirstName,
		LastName:    registerUser.LastName,
		Role:        registerUser.Role,
		Barangay_ID: &barangay_ID_uint,
		Contact:     registerUser.Contact,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (s *UserService) LoginUser(loginUser models.LoginUser) (models.UserStruct, error) {

	validateLoginUser(loginUser)

	var user models.UserStruct
	if err := s.db.Model(&models.User{}).
		Select("id, password, role, barangay_id").
		Where("email = ?", loginUser.Email).
		Scan(&user).Error; err != nil {
		return models.UserStruct{}, fmt.Errorf("%w: email not found", ErrUserNotFound)
	}

	if !CheckPassword(user.Password, loginUser.Password) {
		return models.UserStruct{}, ErrInvalidCredentials
	}

	var barangay models.Barangay
	if err := s.db.Select("name").Where("id = ?", user.Barangay_ID).First(&barangay).Error; err != nil {
		return models.UserStruct{}, fmt.Errorf("failed to retrieve barangay: %w", err)
	}

	user.Barangay_Name = barangay.Name

	return user, nil
}

func (s *UserService) GetUserProfile(userID uint) (models.UserProfile, error) {

	var userProfile models.UserProfile
	if err := s.db.Model(&models.User{}).
		Select("id, email, first_name, last_name, role, contact").
		Where("id = ?", userID).
		Scan(&userProfile).Error; err != nil {
		return models.UserProfile{}, fmt.Errorf("%w: user ID %d not found", ErrUserNotFound, userID)
	}

	return userProfile, nil
}
