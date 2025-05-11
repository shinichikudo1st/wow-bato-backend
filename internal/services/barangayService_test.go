package services

import (
	"database/sql"
	"testing"
	"wow-bato-backend/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestBarangayService_AddNewBarangay(t *testing.T) {
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
	svc := NewBarangayService(gormDB)
	barangay := models.AddBarangay{Name: "Test", City: "TestCity", Region: "TestRegion"}

	// Setup expectations - match what GORM is actually doing
	mock.ExpectBegin()
	// Use AnyArg() for timestamps and other auto-generated fields
	mock.ExpectQuery(`INSERT INTO "barangays"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			"Test", "TestCity", "TestRegion", ""). // Match the actual params
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	// Call the method
	err = svc.AddNewBarangay(barangay)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestBarangayService_DeleteBarangay(t *testing.T) {
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
	svc := NewBarangayService(gormDB)
	barangayID := "1" // Test with a valid ID

	// Setup expectations - match what GORM is actually doing
	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "barangays" WHERE id = \$1`).
		WithArgs(1).                              // The numeric ID after conversion
		WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected
	mock.ExpectCommit()

	// Call the method
	err = svc.DeleteBarangay(barangayID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestBarangayService_UpdateBarangay(t *testing.T) {
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
	svc := NewBarangayService(gormDB)
	barangayID := "1" // Test with a valid ID
	updateData := models.UpdateBarangay{
		Name:   "Updated Name",
		City:   "Updated City",
		Region: "Updated Region",
	}

	// Setup expectations for the first query (finding the barangay)
	rows := sqlmock.NewRows([]string{"id", "name", "city", "region"}).
		AddRow(1, "Old Name", "Old City", "Old Region")

	mock.ExpectQuery(`SELECT (.+) FROM "barangays" WHERE id = \$1`).
		WithArgs(1).
		WillReturnRows(rows)

	// Setup expectations for the update
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "barangays" SET (.+) WHERE "id" = \$`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "Updated Name", "Updated City", "Updated Region", 1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	// Call the method
	err = svc.UpdateBarangay(barangayID, updateData)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
