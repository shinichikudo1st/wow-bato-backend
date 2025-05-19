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
