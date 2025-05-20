package services

import (
	"database/sql"
	"testing"
	"wow-bato-backend/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestBudgetItemService_AddBudgetItem(t *testing.T) {
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

	svc := NewBudgetItemService(gormDB)

	projectID := "1"
	budgetItem := models.NewBudgetItem{
		Name:             "Test Item",
		Amount_Allocated: 1000.0,
		Description:      "Test Description",
		Status:           "Pending",
	}

	// Setup expectations
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "budget_items"`).
		WithArgs(
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), // Created, updated timestamps
			budgetItem.Name,
			budgetItem.Amount_Allocated,
			budgetItem.Description,
			budgetItem.Status,
			uint(1), // ProjectID as uint
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err = svc.AddBudgetItem(projectID, budgetItem)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestBudgetItemService_GetAllBudgetItem(t *testing.T) {
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

	svc := NewBudgetItemService(gormDB)

	projectID := "1"
	filter := "All"
	page := "1"

	// Setup mock rows
	rows := sqlmock.NewRows([]string{"id", "name", "amount_allocated", "description", "status", "project_id"}).
		AddRow(1, "Item 1", 500.0, "Desc 1", "Pending", 1).
		AddRow(2, "Item 2", 1500.0, "Desc 2", "Approved", 1)

	mock.ExpectQuery(`SELECT \* FROM "budget_items" WHERE project_id = \$1 LIMIT \$2 OFFSET \$3`).
		WithArgs(1, 5, 0).
		WillReturnRows(rows)

	items, err := svc.GetAllBudgetItem(projectID, filter, page)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(items))
	}

	if items[0].Name != "Item 1" || items[1].Name != "Item 2" {
		t.Errorf("Unexpected item names: %s, %s", items[0].Name, items[1].Name)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestBudgetItemService_CountBudgetItem(t *testing.T) {
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

	svc := NewBudgetItemService(gormDB)

	projectID := "1"
	expectedCount := int64(3)

	// Setup mock for the count query
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(expectedCount)
	mock.ExpectQuery(`SELECT count\(\*\) FROM "budget_items" WHERE project_id = \$1`).
		WithArgs(1).
		WillReturnRows(countRows)

	count, err := svc.CountBudgetItem(projectID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if count != expectedCount {
		t.Errorf("Expected count %d, got %d", expectedCount, count)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
