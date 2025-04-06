// authentication_test.go
// Package handlers_test provides unit tests for the handlers package.
// Test the authentication.go functions
package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"wow-bato-backend/internal/models"

	"github.com/gin-gonic/gin"
)

func TestRegisterUser(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Setup router
	r := gin.Default()
	r.POST("/register")

	// Test cases
	tests := []struct {
		name          string
		requestBody   models.RegisterUser
		expectedCode  int
		expectedError bool
		expectedBody  string
	}{
		{
			name: "Success",
			requestBody: models.RegisterUser{
				Email:       "test@example.com",
				Password:    "password123",
				FirstName:   "John",
				LastName:    "Doe",
				Role:        "resident",
				Barangay_ID: "1",
				Contact:     "+63 912 345 6789",
			},
			expectedCode: http.StatusOK,
			expectedBody: `{"message":"User registered successfully"}`,
		},
		{
			name: "Empty Email",
			requestBody: models.RegisterUser{
				Password:    "password123",
				FirstName:   "John",
				LastName:    "Doe",
				Role:        "resident",
				Barangay_ID: "1",
				Contact:     "+63 912 345 6789",
			},
			expectedCode:  http.StatusBadRequest,
			expectedError: true,
		},
		{
			name: "Empty Password",
			requestBody: models.RegisterUser{
				Email:       "test@example.com",
				FirstName:   "John",
				LastName:    "Doe",
				Role:        "resident",
				Barangay_ID: "1",
				Contact:     "+63 912 345 6789",
			},
			expectedCode:  http.StatusBadRequest,
			expectedError: true,
		},
		{
			name: "Invalid Email Format",
			requestBody: models.RegisterUser{
				Email:       "invalid-email",
				Password:    "password123",
				FirstName:   "John",
				LastName:    "Doe",
				Role:        "resident",
				Barangay_ID: "1",
				Contact:     "+63 912 345 6789",
			},
			expectedCode:  http.StatusBadRequest,
			expectedError: true,
		},
		{
			name: "Invalid Role",
			requestBody: models.RegisterUser{
				Email:       "test@example.com",
				Password:    "password123",
				FirstName:   "John",
				LastName:    "Doe",
				Role:        "invalid_role",
				Barangay_ID: "1",
				Contact:     "+63 912 345 6789",
			},
			expectedCode:  http.StatusBadRequest,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert request body to JSON
			jsonValue, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Fatalf("Failed to marshal request body: %v", err)
			}

			// Create request
			req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonValue))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			r.ServeHTTP(w, req)

			// Check status code
			if w.Code != tt.expectedCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedCode, w.Code)
			}

			// Check response body for success case
			if tt.expectedCode == http.StatusOK {
				if w.Body.String() != tt.expectedBody {
					t.Errorf("Expected body %s, got %s", tt.expectedBody, w.Body.String())
				}
			}

			// Check for error response when expected
			if tt.expectedError {
				var response map[string]interface{}
				err = json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				if _, exists := response["error"]; !exists {
					t.Error("Expected error message in response, got none")
				}
			}
		})
	}
}
