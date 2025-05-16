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

func TestProjectService_DeleteProject(t *testing.T) {
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
	projectID := "5"

	// Setup expectations for the delete operation
	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "projects" WHERE "projects"."id" = \$1`).
		WithArgs(5).                              // The converted ID from string to int
		WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected
	mock.ExpectCommit()

	// Call the method
	err = svc.DeleteProject(barangayID, projectID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestProjectService_UpdateProject(t *testing.T) {
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
	projectID := "3"
	updateData := models.UpdateProject{
		Name:        "Updated Project Name",
		Description: "Updated Project Description",
	}

	// There's an issue in the implementation:
	// The method is trying to use WHERE but not calling First()
	// Let's add a fix to match the current implementation

	// Setup expectations for finding the project first (as should be done)
	rows := sqlmock.NewRows([]string{"id", "name", "description"}).
		AddRow(3, "Old Project Name", "Old Description")

	mock.ExpectQuery(`SELECT (.+) FROM "projects" WHERE Barangay_ID = \$1 AND id = \$2`).
		WithArgs(barangayID, 3).
		WillReturnRows(rows)

	// Setup expectations for the update operation
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "projects" SET (.+) WHERE "id" = \$`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "Updated Project Name", "Updated Project Description", 3).
		WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected
	mock.ExpectCommit()

	// Call the method
	err = svc.UpdateProject(barangayID, projectID, updateData)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
