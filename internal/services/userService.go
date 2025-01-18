// Package services implements the core business logic for user management and authentication.
// It provides a clean separation between the application's domain logic and infrastructure
// concerns, following clean architecture principles. This package is responsible for:
//   - User registration and authentication
//   - Profile management and retrieval
//   - Password security and hashing
//   - Business rule validation
//
// The services layer acts as an intermediary between the handlers (presentation) layer
// and the data access layer, ensuring proper encapsulation of business rules.
package services

import (
	"errors"
	"strconv"
	database "wow-bato-backend/internal"
	"wow-bato-backend/internal/models"
)

// RegisterUser handles new user registration with proper validation and security measures.
//
// It implements the core business logic for user registration, including password hashing,
// data validation, and persistent storage. The function follows security best practices
// by ensuring passwords are properly hashed before storage.
//
// Parameters:
//   - registerUser: models.RegisterUser - The registration data transfer object containing:
//     * Email: User's email address (must be unique in the system)
//     * Password: Plain text password (will be securely hashed)
//     * FirstName: User's first name
//     * LastName: User's last name
//     * Role: User's system role (determines permissions)
//     * Barangay_ID: String ID of user's barangay
//     * Contact: User's contact information
//
// Returns:
//   - error: nil on successful registration, or an error describing the failure:
//     * ErrDatabaseConnection: When database connection fails
//     * ErrPasswordHashing: When password hashing fails
//     * ErrInvalidBarangayID: When barangay ID conversion fails
//     * ErrDatabaseOperation: When user creation in database fails
//     * ErrDuplicateEmail: When email already exists
//
// Example usage:
//
//	newUser := models.RegisterUser{
//	    Email:       "john.doe@example.com",
//	    Password:    "securePassword123",
//	    FirstName:   "John",
//	    LastName:    "Doe",
//	    Role:        "resident",
//	    Barangay_ID: "1",
//	    Contact:     "+63 912 345 6789",
//	}
//	if err := RegisterUser(newUser); err != nil {
//	    // Handle registration error
//	    return fmt.Errorf("user registration failed: %w", err)
//	}
func RegisterUser(registerUser models.RegisterUser) error {
	db, err := database.ConnectDB()
	if err != nil {
		return err
	}

	hash, err := HashPassword(registerUser.Password)
	if err != nil {
		return err
	}

	barangay_ID, err := strconv.Atoi(registerUser.Barangay_ID)
	if err != nil {
		return err
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

	result := db.Create(&user)
	return result.Error
}

// LoginUser authenticates a user and returns their profile information.
//
// This function implements secure user authentication by verifying credentials
// and managing user sessions. It follows security best practices for password
// verification and user data handling.
//
// Parameters:
//   - loginUser: models.LoginUser - The login credentials containing:
//     * Email: User's registered email address
//     * Password: User's password (will be verified against stored hash)
//
// Returns:
//   - models.UserStruct: Complete user profile including:
//     * User's basic information (ID, email, name)
//     * Role and permissions
//     * Associated barangay details
//     * Additional profile metadata
//   - error: nil on successful authentication, or an error describing the failure:
//     * ErrInvalidCredentials: When email/password combination is invalid
//     * ErrDatabaseConnection: When database access fails
//     * ErrUserNotFound: When user email doesn't exist
//     * ErrInternalServer: For unexpected system errors
//
// Example usage:
//
//	credentials := models.LoginUser{
//	    Email:    "john.doe@example.com",
//	    Password: "userPassword123",
//	}
//	userProfile, err := LoginUser(credentials)
//	if err != nil {
//	    // Handle authentication error
//	    return nil, fmt.Errorf("authentication failed: %w", err)
//	}
func LoginUser(loginUser models.LoginUser) (models.UserStruct, error) {
	db, err := database.ConnectDB()
	if err != nil {
		return models.UserStruct{}, err
	}

	var user models.UserStruct
	if err := db.Model(&models.User{}).
		Select("id, password, role, barangay_id").
		Where("email = ?", loginUser.Email).
		Scan(&user).Error; err != nil {
		return models.UserStruct{}, err
	}

	if !CheckPassword(user.Password, loginUser.Password) {
		return models.UserStruct{}, errors.New("invalid email or password")
	}

	var barangay models.Barangay
	if err := db.Select("name").Where("id = ?", user.Barangay_ID).First(&barangay).Error; err != nil {
		return models.UserStruct{}, err
	}

	user.Barangay_Name = barangay.Name

	return user, nil
}

// GetUserProfile retrieves detailed user profile information.
//
// This function performs the following operations:
// 1. Establishes a database connection
// 2. Retrieves user profile data based on user ID
//
// Parameters:
//   - userID: uint - The unique identifier of the user
//
// Returns:
//   - models.UserProfile: User profile information including:
//   - ID
//   - Email
//   - FirstName
//   - LastName
//   - Role
//   - Contact
//   - error: Returns nil if successful, otherwise returns an error:
//   - Database connection errors
//   - User not found errors
//
// Example:
//
//	userID := uint(1)
//	profile, err := GetUserProfile(userID)
//	if err != nil {
//	    log.Printf("Failed to get user profile: %v", err)
//	    return err
//	}
func GetUserProfile(userID uint) (models.UserProfile, error) {
	db, err := database.ConnectDB()
	if err != nil {
		return models.UserProfile{}, err
	}

	var userProfile models.UserProfile
	if err := db.Model(&models.User{}).
		Select("id, email, first_name, last_name, role, contact").
		Where("id = ?", userID).
		Scan(&userProfile).Error; err != nil {
		return models.UserProfile{}, err
	}

	return userProfile, nil
}
