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

func TestCreateFeedBack(t *testing.T) {
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
	svc := services.NewFeedbackService(gormDB)
	handlersObj := handlers.NewFeedbackHandlers(svc)

	r.POST("/feedback/:projectID", func(c *gin.Context) {
		sess := sessions.Default(c)
		sess.Set("authenticated", true)
		sess.Set("user_id", uint(1))
		sess.Set("user_role", "resident")
		sess.Save()
		handlersObj.CreateFeedBack(c)
	})

	// Setup DB expectations for a successful insert
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "feedbacks"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			"Test feedback", uint(1), "resident", uint(2)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	feedbackBody := models.NewFeedback{
		Content: "Test feedback",
	}
	jsonValue, _ := json.Marshal(feedbackBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/feedback/2", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("New feedback created")) {
		t.Errorf("Expected success message, got %s", w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestGetAllFeedbacks(t *testing.T) {
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
	svc := services.NewFeedbackService(gormDB)
	handlersObj := handlers.NewFeedbackHandlers(svc)

	r.GET("/feedbacks/:projectID", func(c *gin.Context) {
		sess := sessions.Default(c)
		sess.Set("authenticated", true)
		sess.Save()
		handlersObj.GetAllFeedbacks(c)
	})

	// Setup DB expectations for feedbacks and users
	mock.ExpectQuery(`SELECT id, content, role, project_id, user_id FROM "feedbacks" WHERE project_id = \$1`).
		WithArgs(2).
		WillReturnRows(sqlmock.NewRows([]string{"id", "content", "role", "project_id", "user_id"}).
			AddRow(1, "Feedback 1", "resident", 2, 10).
			AddRow(2, "Feedback 2", "admin", 2, 20))
	mock.ExpectQuery(`SELECT id, first_name, last_name FROM "users" WHERE id IN \(\$1,\$2\)`).
		WithArgs(10, 20).
		WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "last_name"}).
			AddRow(10, "John", "Doe").
			AddRow(20, "Jane", "Smith"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/feedbacks/2", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("feedbacks")) {
		t.Errorf("Expected feedbacks in response, got %s", w.Body.String())
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("Feedback 1")) || !bytes.Contains(w.Body.Bytes(), []byte("Feedback 2")) {
		t.Errorf("Expected feedback contents in response, got %s", w.Body.String())
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("John")) || !bytes.Contains(w.Body.Bytes(), []byte("Jane")) {
		t.Errorf("Expected user names in response, got %s", w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
