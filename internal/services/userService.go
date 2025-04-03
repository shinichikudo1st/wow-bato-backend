package services

import (
	"errors"
	"fmt"
	"strconv"
	database "wow-bato-backend/internal"
	"wow-bato-backend/internal/models"
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

func RegisterUser(registerUser models.RegisterUser) error {
	if err := validateUserRegistration(registerUser); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	db, err := database.ConnectDB()
	if err != nil {
		return fmt.Errorf("database connection failed: %w", err)
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

	if result := db.Create(&user); result.Error != nil {
		return fmt.Errorf("failed to create user: %w", result.Error)
	}

	return nil
}

func LoginUser(loginUser models.LoginUser) (models.UserStruct, error) {
	if loginUser.Email == "" {
		return models.UserStruct{}, ErrEmptyEmail
	}
	if loginUser.Password == "" {
		return models.UserStruct{}, ErrEmptyPassword
	}

	db, err := database.ConnectDB()
	if err != nil {
		return models.UserStruct{}, fmt.Errorf("database connection failed: %w", err)
	}

	var user models.UserStruct
	if err := db.Model(&models.User{}).
		Select("id, password, role, barangay_id").
		Where("email = ?", loginUser.Email).
		Scan(&user).Error; err != nil {
		return models.UserStruct{}, fmt.Errorf("%w: email not found", ErrUserNotFound)
	}

	if !CheckPassword(user.Password, loginUser.Password) {
		return models.UserStruct{}, ErrInvalidCredentials
	}

	var barangay models.Barangay
	if err := db.Select("name").Where("id = ?", user.Barangay_ID).First(&barangay).Error; err != nil {
		return models.UserStruct{}, fmt.Errorf("failed to retrieve barangay: %w", err)
	}

	user.Barangay_Name = barangay.Name

	return user, nil
}

func GetUserProfile(userID uint) (models.UserProfile, error) {
	db, err := database.ConnectDB()
	if err != nil {
		return models.UserProfile{}, fmt.Errorf("database connection failed: %w", err)
	}

	var userProfile models.UserProfile
	if err := db.Model(&models.User{}).
		Select("id, email, first_name, last_name, role, contact").
		Where("id = ?", userID).
		Scan(&userProfile).Error; err != nil {
		return models.UserProfile{}, fmt.Errorf("%w: user ID %d not found", ErrUserNotFound, userID)
	}

	return userProfile, nil
}
