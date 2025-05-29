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

func TestAddBarangay(t *testing.T) {
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
	svc := services.NewBarangayService(gormDB)
	handlersObj := handlers.NewBarangayHandlers(svc)

	r.POST("/barangay", func(c *gin.Context) {
		sess := sessions.Default(c)
		sess.Set("authenticated", true)
		sess.Save()
		handlersObj.AddBarangay(c)
	})

	// Setup DB expectations for a successful insert
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "barangays"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			"TestBarangay", "TestCity", "TestRegion").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	reqBody := models.AddBarangay{
		Name:   "TestBarangay",
		City:   "TestCity",
		Region: "TestRegion",
	}
	jsonValue, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/barangay", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("Successfully Added New Barangay")) {
		t.Errorf("Expected success message, got %s", w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestGetAllBarangay(t *testing.T) {
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
	svc := services.NewBarangayService(gormDB)
	handlersObj := handlers.NewBarangayHandlers(svc)

	r.GET("/barangay", func(c *gin.Context) {
		sess := sessions.Default(c)
		sess.Set("authenticated", true)
		sess.Save()
		handlersObj.GetAllBarangay(c)
	})

	// Setup DB expectations for a successful select
	mock.ExpectQuery(`SELECT id, name, city, region FROM "barangays" LIMIT \$1 OFFSET \$2`).
		WithArgs(10, 0).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "city", "region"}).
			AddRow(1, "Barangay1", "City1", "Region1").
			AddRow(2, "Barangay2", "City2", "Region2"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/barangay?limit=10&page=1", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("Successfully fetched Barangays")) {
		t.Errorf("Expected success message, got %s", w.Body.String())
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("Barangay1")) || !bytes.Contains(w.Body.Bytes(), []byte("Barangay2")) {
		t.Errorf("Expected barangay names in response, got %s", w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestGetSingleBarangay(t *testing.T) {
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
	svc := services.NewBarangayService(gormDB)
	handlersObj := handlers.NewBarangayHandlers(svc)

	r.GET("/barangay/:barangay_ID", func(c *gin.Context) {
		sess := sessions.Default(c)
		sess.Set("authenticated", true)
		sess.Save()
		handlersObj.GetSingleBarangay(c)
	})

	// Setup DB expectations for a successful select
	mock.ExpectQuery(`SELECT id, name, city, region FROM "barangays" WHERE ID = \$1 LIMIT 1`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "city", "region"}).
			AddRow(1, "Barangay1", "City1", "Region1"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/barangay/1", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("Retrieved specific barangay")) {
		t.Errorf("Expected success message, got %s", w.Body.String())
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("Barangay1")) {
		t.Errorf("Expected barangay name in response, got %s", w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
