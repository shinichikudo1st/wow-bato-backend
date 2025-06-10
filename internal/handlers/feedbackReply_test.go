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
