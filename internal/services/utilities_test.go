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

func TestCheckPassword(t *testing.T) {
	// Create some test passwords and their hashes
	testCases := []struct {
		name        string
		password    string
		inputPass   string
		wantMatch   bool
		setupHash   bool // if true, we'll create a hash first
	}{
		{
			name:      "Correct Password",
			password:  "mySecurePassword123",
			inputPass: "mySecurePassword123",
			wantMatch: true,
			setupHash: true,
		},
		{
			name:      "Incorrect Password",
			password:  "mySecurePassword123",
			inputPass: "wrongPassword123",
			wantMatch: false,
			setupHash: true,
		},
		{
			name:      "Empty Password",
			password:  "mySecurePassword123",
			inputPass: "",
			wantMatch: false,
			setupHash: true,
		},
		{
			name:      "Empty Hash",
			password:  "mySecurePassword123",
			inputPass: "mySecurePassword123",
			wantMatch: false,
			setupHash: false, // don't create hash, use empty string
		},
		{
			name:      "Case Sensitive Check",
			password:  "mySecurePassword123",
			inputPass: "MySecurePassword123",
			wantMatch: false,
			setupHash: true,
		},
		{
			name:      "With Special Characters",
			password:  "my!@#$%^&*()Pass",
			inputPass: "my!@#$%^&*()Pass",
			wantMatch: true,
			setupHash: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var hashedPassword string
			if tt.setupHash {
				var err error
				hashedPassword, err = HashPassword(tt.password)
				if err != nil {
					t.Fatalf("Failed to create hash for test: %v", err)
				}
			}

			// Test password verification
			got := CheckPassword(hashedPassword, tt.inputPass)
			if got != tt.wantMatch {
				t.Errorf("CheckPassword() = %v, want %v", got, tt.wantMatch)
			}

			// Additional security checks for matching passwords
			if tt.wantMatch && tt.setupHash {
				// Verify that the function is consistent
				secondCheck := CheckPassword(hashedPassword, tt.inputPass)
				if !secondCheck {
					t.Error("CheckPassword() not consistent between calls")
				}

				// Verify that adding extra characters fails
				extraChar := CheckPassword(hashedPassword, tt.inputPass+"extra")
				if extraChar {
					t.Error("CheckPassword() accepted password with extra characters")
				}
			}
		})
	}
}
