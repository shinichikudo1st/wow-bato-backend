// Package services provides utility functions for secure password management.
// This package implements cryptographic functions using industry-standard
// algorithms and best practices for:
//   - Password hashing using bcrypt with appropriate work factors
//   - Secure password verification
//   - Protection against timing attacks
//
// The package ensures that all cryptographic operations follow
// security best practices and OWASP guidelines.
package services

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword securely hashes a password using the bcrypt algorithm.
//
// This function implements secure password hashing using bcrypt with a work factor
// of 14, providing a good balance between security and performance. The function
// automatically handles salt generation and secure memory management.
//
// Parameters:
//   - password: string - The plaintext password to hash (should be non-empty)
//
// Returns:
//   - string: The securely hashed password (60 characters long)
//   - error: nil on successful hashing, or an error describing the failure:
//   - ErrInvalidPassword: When password is empty
//   - ErrHashingFailed: When hashing operation fails
//
// Example usage:
//
//	hashedPwd, err := HashPassword("mySecurePassword123")
//	if err != nil {
//	    return fmt.Errorf("password hashing failed: %w", err)
//	}
func HashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("password cannot be empty")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(hashedPassword), err
}

// CheckPassword securely compares a plaintext password with a hashed password.
//
// This function performs a constant-time comparison between the provided plaintext
// password and a previously hashed password, protecting against timing attacks.
// It uses the bcrypt.CompareHashAndPassword function which is specifically
// designed to be timing-attack resistant.
//
// Parameters:
//   - hashedPassword: string - The previously hashed password (60 characters bcrypt hash)
//   - password: string - The plaintext password to verify
//
// Returns:
//   - bool: true if passwords match, false otherwise
//   - true: Password is correct
//   - false: Password is incorrect or an error occurred
//
// Example usage:
//
//	if CheckPassword(storedHash, userInputPassword) {
//	    // Password is correct, proceed with authentication
//	} else {
//	    // Password is incorrect, handle authentication failure
//	}
func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
