package services

import (
	"database/sql"
	"testing"
	"wow-bato-backend/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestBudgetCategoryService_AddBudgetCategory(t *testing.T) {
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
	svc := NewBudgetCategoryService(gormDB)

	// Test data
	budgetCategory := models.NewBudgetCategory{
		Name:        "Test Category",
		Description: "Test Description",
		Barangay_ID: uint(1),
	}

	// Setup expectations - match what GORM is actually doing
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "budget_categories"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			"Test Category", "Test Description", uint(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	// Call the method
	err = svc.AddBudgetCategory(budgetCategory)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
