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

func TestBudgetCategoryService_DeleteBudgetCategory(t *testing.T) {
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
	categoryID := "1"

	// Setup expectations for the delete operation
	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "budget_categories" WHERE id = \$1`).
		WithArgs(1).                              // The converted ID from string to int
		WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected
	mock.ExpectCommit()

	// Call the method
	err = svc.DeleteBudgetCategory(categoryID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestBudgetCategoryService_UpdateBudgetCategory(t *testing.T) {
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
	categoryID := "1"
	updateData := models.UpdateBudgetCategory{
		Name:        "Updated Category",
		Description: "Updated Description",
	}

	// Setup expectations for finding the category first
	findRows := sqlmock.NewRows([]string{"id", "name", "description", "barangay_id"}).
		AddRow(1, "Old Category", "Old Description", 1)

	mock.ExpectQuery(`SELECT (.+) FROM "budget_categories" WHERE id = \$1`).
		WithArgs(1).
		WillReturnRows(findRows)

	// Setup expectations for the update operation
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "budget_categories" SET (.+) WHERE "id" = \$`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "Updated Category", "Updated Description", 1).
		WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected
	mock.ExpectCommit()

	// Call the method
	err = svc.UpdateBudgetCategory(categoryID, updateData)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestBudgetCategoryService_GetAllBudgetCategory(t *testing.T) {
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

	// Test parameters
	barangayID := "1"
	limit := "10"
	page := "1"

	// Setup mock rows with the expected columns from the JOIN query
	rows := sqlmock.NewRows([]string{
		"id", "name", "description", "barangay_id", "project_count",
	}).
		AddRow(1, "Category 1", "Description 1", 1, 2).
		AddRow(2, "Category 2", "Description 2", 1, 0)

	// The query is a complex JOIN with GROUP BY, so use a generic pattern
	mock.ExpectQuery(`SELECT (.+) FROM "budget_categories" LEFT JOIN projects (.+) WHERE (.+) GROUP BY (.+) LIMIT (.+) OFFSET (.+)`).
		WithArgs(1, 10, 0). // barangay_id, limit, offset
		WillReturnRows(rows)

	// Call the method
	results, err := svc.GetAllBudgetCategory(barangayID, limit, page)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Check results length
	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}

	// Verify first result
	if results[0].ID != 1 {
		t.Errorf("Expected ID 1, got %d", results[0].ID)
	}

	if results[0].Name != "Category 1" {
		t.Errorf("Expected name 'Category 1', got %s", results[0].Name)
	}

	if results[0].Description != "Description 1" {
		t.Errorf("Expected description 'Description 1', got %s", results[0].Description)
	}

	if results[0].Barangay_ID != 1 {
		t.Errorf("Expected barangay_id 1, got %d", results[0].Barangay_ID)
	}

	if results[0].ProjectCount != 2 {
		t.Errorf("Expected project_count 2, got %d", results[0].ProjectCount)
	}

	// Verify second result
	if results[1].ID != 2 {
		t.Errorf("Expected ID 2, got %d", results[1].ID)
	}

	if results[1].ProjectCount != 0 {
		t.Errorf("Expected project_count 0, got %d", results[1].ProjectCount)
	}

	// Verify all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
