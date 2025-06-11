package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"wow-bato-backend/internal/handlers"
	"wow-bato-backend/internal/models"
	"wow-bato-backend/internal/services"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestCreateFeedbackReply(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	db, mock, err := sqlmock.New()
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
	svc := services.NewFeedbackReplyService(gormDB)
	handlersObj := handlers.NewFeedbackReplyHandlers(svc)

	r.POST("/feedback-reply/create/:feedbackID", func(c *gin.Context) {
		sess := sessions.Default(c)
		sess.Set("user_id", uint(1))
		sess.Save()
		handlersObj.CreateFeedbackReply(c)
	})

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "feedback_replies"`).
		WithArgs(
			sqlmock.AnyArg(), // gorm.Model fields (ID, CreatedAt, etc.)
			"Test reply",     // content
			2,                // feedback_id
			1,                // user_id
		).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	reqBody := models.Reply{
		Content: "Test reply",
	}
	jsonValue, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/feedback-reply/create/2", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("Reply submitted")) {
		t.Errorf("Expected success message, got %s", w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestGetAllReplies(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.Default()

	db, mock, err := sqlmock.New()
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
	svc := services.NewFeedbackReplyService(gormDB)
	handlersObj := handlers.NewFeedbackReplyHandlers(svc)

	r.GET("/feedback-reply/all/:feedbackID", func(c *gin.Context) {
		handlersObj.GetAllReplies(c)
	})

	mock.ExpectQuery(`SELECT \* FROM "feedback_replies" WHERE feedback_ID = \$1`).
		WithArgs(2).
		WillReturnRows(sqlmock.NewRows([]string{"id", "content", "feedback_id", "user_id"}).
			AddRow(1, "Reply 1", 2, 1).
			AddRow(2, "Reply 2", 2, 2))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/feedback-reply/all/2", nil)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("Reply 1")) || !bytes.Contains(w.Body.Bytes(), []byte("Reply 2")) {
		t.Errorf("Expected reply contents in response, got %s", w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestDeleteFeedbackReply(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.Default()

	db, mock, err := sqlmock.New()
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
	svc := services.NewFeedbackReplyService(gormDB)
	handlersObj := handlers.NewFeedbackReplyHandlers(svc)

	r.DELETE("/feedback-reply/delete/:feedbackID", func(c *gin.Context) {
		handlersObj.DeleteFeedbackReply(c)
	})

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "feedback_replies" WHERE "feedback_replies"."id" = \$1`).
		WithArgs(2).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/feedback-reply/delete/2", nil)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("Reply deleted")) {
		t.Errorf("Expected success message, got %s", w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestEditFeedbackReply(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	db, mock, err := sqlmock.New()
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
	svc := services.NewFeedbackReplyService(gormDB)
	handlersObj := handlers.NewFeedbackReplyHandlers(svc)

	r.PUT("/feedback-reply/edit/:replyID", func(c *gin.Context) {
		sess := sessions.Default(c)
		sess.Set("user_id", uint(1))
		sess.Save()
		handlersObj.EditFeedbackReply(c)
	})

	// Mock the SELECT for existence check
	mock.ExpectQuery(`SELECT \* FROM "feedback_replies" WHERE id = \$1 ORDER BY "feedback_replies"."id" LIMIT 1`).
		WithArgs(2).
		WillReturnRows(sqlmock.NewRows([]string{"id", "content", "feedback_id", "user_id"}).
			AddRow(2, "Old reply", 2, 1))
	// Mock the UPDATE
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "feedback_replies" SET (.+) WHERE "id" = \$1`).
		WithArgs("Updated reply", 2).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	editReply := models.EditReply{
		Content: "Updated reply",
		UserID:  1,
	}
	jsonValue, _ := json.Marshal(editReply)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/feedback-reply/edit/2", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("Reply Edited")) {
		t.Errorf("Expected success message, got %s", w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
