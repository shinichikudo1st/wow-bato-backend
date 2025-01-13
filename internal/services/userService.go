// Package services provides user-related business logic and operations for the application.
// It handles user authentication, registration, and profile management while
// maintaining separation of concerns from the database and presentation layers.
package services

import (
	"errors"
	"strconv"
	database "wow-bato-backend/internal"
	"wow-bato-backend/internal/models"
)

// RegisterUser handles the user registration process in the system.
//
// This function performs the following operations:
// 1. Establishes a database connection
// 2. Hashes the user's password for secure storage
// 3. Converts the barangay ID from string to uint
// 4. Creates a new user record in the database
//
// Parameters:
//   - registerUser: models.RegisterUser - Contains user registration data including:
//   - Email: User's email address (must be unique)
//   - Password: User's plain text password (will be hashed)
//   - FirstName: User's first name
//   - LastName: User's last name
//   - Role: User's role in the system
//   - Barangay_ID: String ID of user's barangay (will be converted to uint)
//   - Contact: User's contact information
//
// Returns:
//   - error: Returns nil if registration is successful, otherwise returns an error:
//   - Database connection errors
//   - Password hashing errors
//   - Barangay ID conversion errors
//   - Database creation errors
//
// Example:
//
//	user := models.RegisterUser{
//	    Email: "user@example.com",
//	    Password: "securepassword123",
//	    FirstName: "John",
//	    LastName: "Doe",
//	    Role: "user",
//	    Barangay_ID: "1",
//	    Contact: "1234567890",
//	}
//	err := RegisterUser(user)
//	if err != nil {
//	    log.Printf("Registration failed: %v", err)
//	    return err
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

// LoginUser authenticates a user and returns their information.
//
// This function performs the following operations:
// 1. Establishes a database connection
// 2. Retrieves user information based on email
// 3. Verifies the provided password
// 4. Fetches and includes the barangay name in the response
//
// Parameters:
//   - loginUser: models.LoginUser - Contains login credentials:
//   - Email: User's email address
//   - Password: User's password (will be verified against stored hash)
//
// Returns:
//   - models.UserStruct: User information including ID, role, and barangay details
//   - error: Returns nil if login is successful, otherwise returns an error:
//   - Database connection errors
//   - User not found
//   - Invalid password
//   - Barangay lookup errors
//
// Example:
//
//	credentials := models.LoginUser{
//	    Email: "user@example.com",
//	    Password: "userpassword123",
//	}
//	userInfo, err := LoginUser(credentials)
//	if err != nil {
//	    log.Printf("Login failed: %v", err)
//	    return err
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
