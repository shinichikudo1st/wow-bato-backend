// authentication_test.go
// Package handlers_test provides unit tests for the handlers package.
// Test the authentication.go functions
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

// MockUserService implements the RegisterUser method for testing
// You can expand this mock for more methods as needed

type MockUserService struct {
	RegisterUserFunc func(models.RegisterUser) error
}

func (m *MockUserService) RegisterUser(u models.RegisterUser) error {
	if m.RegisterUserFunc != nil {
		return m.RegisterUserFunc(u)
	}
	return nil
}

func TestRegisterUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// Setup mock DB
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

	svc := services.NewUserService(gormDB)
	handlersObj := handlers.NewUserHandlers(svc)

	r := gin.Default()
	r.POST("/register", handlersObj.RegisterUser)

	// Setup DB expectations for a successful registration
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "users"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			"test@example.com", sqlmock.AnyArg(), "John", "Doe", "resident", uint(1), "+63 912 345 6789").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	reqBody := models.RegisterUser{
		Email:       "test@example.com",
		Password:    "password123",
		FirstName:   "John",
		LastName:    "Doe",
		Role:        "resident",
		Barangay_ID: "1",
		Contact:     "+63 912 345 6789",
	}
	jsonValue, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("User registered successfully")) {
		t.Errorf("Expected success message, got %s", w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestLoginUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// Setup mock DB
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

	svc := services.NewUserService(gormDB)
	handlersObj := handlers.NewUserHandlers(svc)

	r := gin.Default()
	r.POST("/login", handlersObj.LoginUser)

	// Setup DB expectations for a successful login
	hashedPassword, _ := services.HashPassword("password123")
	mock.ExpectQuery(`SELECT id, password, role, barangay_id FROM "users" WHERE email = \$1`).
		WithArgs("test@example.com").
		WillReturnRows(sqlmock.NewRows([]string{"id", "password", "role", "barangay_id"}).
			AddRow(1, hashedPassword, "resident", 1))
	mock.ExpectQuery(`SELECT name FROM "barangays" WHERE id = \$1`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("Barangay Uno"))

	// Prepare request body
	loginReq := models.LoginUser{
		Email:    "test@example.com",
		Password: "password123",
	}
	jsonValue, _ := json.Marshal(loginReq)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("User logged in successfully")) {
		t.Errorf("Expected success message, got %s", w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestLogoutUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	db, _, err := sqlmock.New()
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
	svc := services.NewUserService(gormDB)
	handlersObj := handlers.NewUserHandlers(svc)

	r.POST("/logout", handlersObj.LogoutUser)

	// Simulate a logged-in session
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/logout", nil)
	// Set a session cookie (Gin will create a new session if not present, so this is enough for this test)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("User logged out successfully")) {
		t.Errorf("Expected logout message, got %s", w.Body.String())
	}
}

func TestCheckAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	db, _, err := sqlmock.New()
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
	svc := services.NewUserService(gormDB)
	handlersObj := handlers.NewUserHandlers(svc)

	r.GET("/checkAuth", func(c *gin.Context) {
		sess := sessions.Default(c)
		sess.Set("authenticated", true)
		sess.Set("user_role", "resident")
		sess.Set("user_id", 1)
		sess.Set("barangay_id", 1)
		sess.Set("barangay_name", "Barangay Uno")
		sess.Save()
		handlersObj.CheckAuth(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/checkAuth", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("sessionStatus")) {
		t.Errorf("Expected sessionStatus in response, got %s", w.Body.String())
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("resident")) {
		t.Errorf("Expected role in response, got %s", w.Body.String())
	}
}

func TestGetUserProfile(t *testing.T) {
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
	svc := services.NewUserService(gormDB)
	handlersObj := handlers.NewUserHandlers(svc)

	// Mock user profile data retrieval
	profileRows := sqlmock.NewRows([]string{"id", "email", "first_name", "last_name", "role", "contact"}).
		AddRow(1, "test@example.com", "John", "Doe", "resident", "+63 912 345 6789")
	mock.ExpectQuery(`SELECT id, email, first_name, last_name, role, contact FROM "users" WHERE id = \$1`).
		WithArgs(1).
		WillReturnRows(profileRows)

	r.GET("/profile", func(c *gin.Context) {
		sess := sessions.Default(c)
		sess.Set("authenticated", true)
		sess.Set("user_id", 1)
		sess.Save()
		handlersObj.GetUserProfile(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/profile", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("User profile fetched successfully")) {
		t.Errorf("Expected profile message, got %s", w.Body.String())
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("test@example.com")) {
		t.Errorf("Expected email in response, got %s", w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
