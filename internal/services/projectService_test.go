package services

import (
	"database/sql"
	"testing"
	"time"
	"wow-bato-backend/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestProjectService_AddNewProject(t *testing.T) {
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
	svc := NewProjectService(gormDB)

	// Test parameters
	barangayID := uint(1)
	categoryID := "2"

	// Test project data
	newProject := models.NewProject{
		Name:        "Test Project",
		Description: "Test Project Description",
		StartDate:   "2023-01-01",
		EndDate:     "2023-12-31",
		Status:      "planned",
	}

	// Parse the dates for verification
	startDate, _ := time.Parse(GO_DATE_FORMAT, newProject.StartDate)
	endDate, _ := time.Parse(GO_DATE_FORMAT, newProject.EndDate)

	// Setup expectations - match what GORM is actually doing
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "projects"`).
		WithArgs(
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), // Created, updated timestamps
			barangayID, uint(2), // Barangay ID and category ID
			newProject.Name, newProject.Description,
			startDate, endDate, newProject.Status,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	// Call the method
	err = svc.AddNewProject(barangayID, categoryID, newProject)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
