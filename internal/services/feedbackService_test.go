package services

import (
	"database/sql"
	"testing"
	"wow-bato-backend/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestFeedbackService_CreateFeedback(t *testing.T) {
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
	svc := NewFeedbackService(gormDB)

	// Test feedback data
	newFeedback := models.CreateFeedback{
		Content:   "Test feedback content",
		UserID:    1,
		Role:      "user",
		ProjectID: 1,
	}

	// Setup expectations
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "feedbacks"`).
		WithArgs(
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), // Created, updated timestamps
			newFeedback.Content,
			newFeedback.UserID,
			newFeedback.Role,
			newFeedback.ProjectID,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	// Call the method
	err = svc.CreateFeedback(newFeedback)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
