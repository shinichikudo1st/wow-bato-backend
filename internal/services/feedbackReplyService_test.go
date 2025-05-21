package services

import (
	"database/sql"
	"testing"
	"wow-bato-backend/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestFeedbackReplyService_CreateFeedbackReply(t *testing.T) {
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

	svc := NewFeedbackReplyService(gormDB)

	newReply := models.NewFeedbackReply{
		Content:    "Test reply content",
		FeedbackID: "1",
		UserID:     2,
	}

	// Setup expectations
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "feedback_replies"`).
		WithArgs(
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), // Created, updated timestamps
			newReply.Content,
			uint(1), // FeedbackID as uint
			newReply.UserID,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err = svc.CreateFeedbackReply(newReply)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestFeedbackReplyService_GetAllReplies(t *testing.T) {
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

	svc := NewFeedbackReplyService(gormDB)

	feedbackID := "1"

	// Setup mock rows
	rows := sqlmock.NewRows([]string{"id", "content", "feedback_id", "user_id"}).
		AddRow(1, "Reply 1", 1, 2).
		AddRow(2, "Reply 2", 1, 3)

	mock.ExpectQuery(`SELECT \* FROM "feedback_replies" WHERE feedback_ID = \$1`).
		WithArgs(1).
		WillReturnRows(rows)

	replies, err := svc.GetAllReplies(feedbackID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(replies) != 2 {
		t.Errorf("Expected 2 replies, got %d", len(replies))
	}

	if replies[0].Content != "Reply 1" || replies[1].Content != "Reply 2" {
		t.Errorf("Unexpected reply contents: %s, %s", replies[0].Content, replies[1].Content)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestFeedbackReplyService_DeleteFeedbackReply(t *testing.T) {
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

	svc := NewFeedbackReplyService(gormDB)

	replyID := "3"

	// Mock the delete operation
	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "feedback_replies" WHERE id = \$1`).
		WithArgs(3).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err = svc.DeleteFeedbackReply(replyID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
