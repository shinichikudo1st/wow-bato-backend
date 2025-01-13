package services

import "golang.org/x/crypto/bcrypt"

// HashPassword hashes the user's password for secure storage
//
// This function performs the following operations:
// 1. Hashes the user's password using bcrypt
// 2. Returns the hashed password and any errors
//
// Parameters:
//   - password: string - The user's plain text password
//
// Returns:
//   - string: The hashed password
//   - error: Returns nil if successful, otherwise returns an error
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(hashedPassword), err
}

// CheckPassword checks if the provided password matches the hashed password
//
// This function performs the following operations:
// 1. Compares the provided password with the hashed password
// 2. Returns true if the passwords match, false otherwise
//
// Parameters:
//   - hashedPassword: string - The hashed password from the database
//   - password: string - The user's plain text password
//
// Returns:
//   - bool: Returns true if the passwords match, false otherwise
func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
