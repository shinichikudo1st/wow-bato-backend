package services

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name        string
		password    string
		wantErr     bool
		checkLength bool // bcrypt hashes should be 60 characters
	}{
		{
			name:        "Valid Password",
			password:    "mySecurePassword123",
			wantErr:     false,
			checkLength: true,
		},
		{
			name:        "Empty Password",
			password:    "",
			wantErr:     true,
			checkLength: false,
		},
		{
			name:        "Long Password",
			password:    "ThisIsAVeryLongPasswordThatShouldStillWorkFineWithBcrypt123!@#",
			wantErr:     false,
			checkLength: true,
		},
		{
			name:        "Special Characters",
			password:    "!@#$%^&*()_+-=[]{}|;:,.<>?",
			wantErr:     false,
			checkLength: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test hashing
			hash, err := HashPassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Check hash length (bcrypt hashes are always 60 characters)
				if tt.checkLength && len(hash) != 60 {
					t.Errorf("HashPassword() hash length = %v, want 60", len(hash))
				}

				// Verify the hash is valid by comparing with original password
				err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(tt.password))
				if err != nil {
					t.Errorf("HashPassword() produced invalid hash, verification failed: %v", err)
				}

				// Verify the hash is different from the original password (basic security check)
				if hash == tt.password {
					t.Error("HashPassword() returned plaintext password instead of hash")
				}

				// Verify that two hashes of the same password are different (salt test)
				hash2, _ := HashPassword(tt.password)
				if hash == hash2 {
					t.Error("HashPassword() produced identical hashes for same password")
				}
			}
		})
	}
}
