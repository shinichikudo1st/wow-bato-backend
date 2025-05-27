package services

import (
	"database/sql"
	"testing"
	"wow-bato-backend/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestUserService_RegisterUser(t *testing.T) {
	
	var db *sql.DB
	var mock sqlmock.Sqlmock
	var err error

	db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	dialector := postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open gorm: %v", err)
	}

	svc := NewUserService(gormDB)
	registerUser := models.RegisterUser{
		Email:       "test@example.com",
		Password:    "password123",
		FirstName:   "Test",
		LastName:    "User",
		Role:        "user",
		Barangay_ID: "1",
		Contact:     "1234567890",
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "users"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			"test@example.com", sqlmock.AnyArg(), "Test", "User", "user", uint(1), "1234567890").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err = svc.RegisterUser(registerUser)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestUserService_LoginUser(t *testing.T) {
	var db *sql.DB
	var mock sqlmock.Sqlmock
	var err error

	db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	dialector := postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open gorm: %v", err)
	}

	svc := NewUserService(gormDB)

	hashedPassword, err := HashPassword("password123")
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	loginUser := models.LoginUser{
		Email:    "test@example.com",
		Password: "password123",
	}

	userRows := sqlmock.NewRows([]string{"id", "password", "role", "barangay_id"}).
		AddRow(1, hashedPassword, "user", 1)

	mock.ExpectQuery(`SELECT id, password, role, barangay_id FROM "users" WHERE email = \$1`).
		WithArgs("test@example.com").
		WillReturnRows(userRows)

	barangayRows := sqlmock.NewRows([]string{"name"}).
		AddRow("Test Barangay")

	mock.ExpectQuery(`SELECT name FROM "barangays" WHERE id = \$1`).
		WithArgs(1).
		WillReturnRows(barangayRows)

	result, err := svc.LoginUser(loginUser)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.ID != 1 {
		t.Errorf("Expected user ID 1, got %d", result.ID)
	}

	if result.Role != "user" {
		t.Errorf("Expected role 'user', got %s", result.Role)
	}

	if result.Barangay_ID != 1 {
		t.Errorf("Expected barangay ID 1, got %d", result.Barangay_ID)
	}

	if result.Barangay_Name != "Test Barangay" {
		t.Errorf("Expected barangay name 'Test Barangay', got %s", result.Barangay_Name)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestUserService_GetUserProfile(t *testing.T) {
	// Create a new SQL mock
	var db *sql.DB
	var mock sqlmock.Sqlmock
	var err error

	db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	// Connect GORM to the mock database
	dialector := postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open gorm: %v", err)
	}

	// Create the service with the mocked database
	svc := NewUserService(gormDB)

	// Test user ID
	userID := uint(1)

	// Mock user profile data retrieval
	profileRows := sqlmock.NewRows([]string{"id", "email", "first_name", "last_name", "role", "contact"}).
		AddRow(1, "test@example.com", "Test", "User", "user", "1234567890")

	mock.ExpectQuery(`SELECT id, email, first_name, last_name, role, contact FROM "users" WHERE id = \$1`).
		WithArgs(userID).
		WillReturnRows(profileRows)

	// Call the method
	result, err := svc.GetUserProfile(userID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify results
	if result.ID != 1 {
		t.Errorf("Expected user ID 1, got %d", result.ID)
	}

	if result.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got %s", result.Email)
	}

	if result.FirstName != "Test" {
		t.Errorf("Expected first name 'Test', got %s", result.FirstName)
	}

	if result.LastName != "User" {
		t.Errorf("Expected last name 'User', got %s", result.LastName)
	}

	if result.Role != "user" {
		t.Errorf("Expected role 'user', got %s", result.Role)
	}

	if result.Contact != "1234567890" {
		t.Errorf("Expected contact '1234567890', got %s", result.Contact)
	}

	// Verify all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
