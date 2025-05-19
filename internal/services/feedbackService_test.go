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

func TestFeedbackService_GetAllFeedback(t *testing.T) {
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

	svc := NewFeedbackService(gormDB)

	projectID := "1"
	projectIDInt := 1

	// Mock feedbacks returned by the first query
	feedbackRows := sqlmock.NewRows([]string{"id", "content", "role", "project_id", "user_id"}).
		AddRow(1, "Feedback 1", "user", projectIDInt, 10).
		AddRow(2, "Feedback 2", "admin", projectIDInt, 20)

	mock.ExpectQuery(`SELECT id, content, role, project_id, user_id FROM "feedbacks" WHERE project_id = \$1`).
		WithArgs(projectIDInt).
		WillReturnRows(feedbackRows)

	// Mock users returned by the second query
	userRows := sqlmock.NewRows([]string{"id", "first_name", "last_name"}).
		AddRow(10, "John", "Doe").
		AddRow(20, "Jane", "Smith")

	mock.ExpectQuery(`SELECT id, first_name, last_name FROM "users" WHERE id IN \(\$1,\$2\)`).
		WithArgs(10, 20).
		WillReturnRows(userRows)

	// Call the method
	feedbacks, err := svc.GetAllFeedback(projectID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(feedbacks) != 2 {
		t.Errorf("Expected 2 feedbacks, got %d", len(feedbacks))
	}

	if feedbacks[0].FirstName != "John" || feedbacks[0].LastName != "Doe" {
		t.Errorf("Expected first feedback user to be John Doe, got %s %s", feedbacks[0].FirstName, feedbacks[0].LastName)
	}
	if feedbacks[1].FirstName != "Jane" || feedbacks[1].LastName != "Smith" {
		t.Errorf("Expected second feedback user to be Jane Smith, got %s %s", feedbacks[1].FirstName, feedbacks[1].LastName)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
