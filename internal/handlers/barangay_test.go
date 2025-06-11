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

func TestDeleteBarangay(t *testing.T) {
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

	r.DELETE("/barangay/:barangay_ID", func(c *gin.Context) {
		sess := sessions.Default(c)
		sess.Set("authenticated", true)
		sess.Save()
		handlersObj.DeleteBarangay(c)
	})

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "barangays" WHERE id = \$1`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/barangay/1", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("Successfully deleted the Barangay")) {
		t.Errorf("Expected success message, got %s", w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestUpdateBarangay(t *testing.T) {
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

	r.PUT("/barangay/:barangay_ID", func(c *gin.Context) {
		sess := sessions.Default(c)
		sess.Set("authenticated", true)
		sess.Save()
		handlersObj.UpdateBarangay(c)
	})

	findRows := sqlmock.NewRows([]string{"id", "name", "city", "region"}).
		AddRow(1, "OldName", "OldCity", "OldRegion")
	mock.ExpectQuery(`SELECT \* FROM "barangays" WHERE id = \$1 ORDER BY "barangays"."id" LIMIT 1`).
		WithArgs(1).
		WillReturnRows(findRows)
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "barangays" SET (.+) WHERE "id" = \$1`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "NewName", "NewCity", "NewRegion", 1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	updateBody := models.UpdateBarangay{
		Name:   "NewName",
		City:   "NewCity",
		Region: "NewRegion",
	}
	jsonValue, _ := json.Marshal(updateBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/barangay/1", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("Successfully Updated Barangay")) {
		t.Errorf("Expected success message, got %s", w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestGetBarangayOptions(t *testing.T) {
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

	r.GET("/barangay/options", func(c *gin.Context) {
		handlersObj.GetBarangayOptions(c)
	})

	mock.ExpectQuery(`SELECT id, name FROM "barangays"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "Barangay1").
			AddRow(2, "Barangay2"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/barangay/options", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("Barangays found")) {
		t.Errorf("Expected success message, got %s", w.Body.String())
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("Barangay1")) || !bytes.Contains(w.Body.Bytes(), []byte("Barangay2")) {
		t.Errorf("Expected barangay names in response, got %s", w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestGetPublicBarangay(t *testing.T) {
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

	r.GET("/barangay/public", func(c *gin.Context) {
		handlersObj.GetPublicBarangay(c)
	})

	mock.ExpectQuery(`SELECT \* FROM "public_barangay_displays"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "city", "region"}).
			AddRow(1, "Barangay1", "City1", "Region1").
			AddRow(2, "Barangay2", "City2", "Region2"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/barangay/public", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("All barangays retrieved")) {
		t.Errorf("Expected success message, got %s", w.Body.String())
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("Barangay1")) || !bytes.Contains(w.Body.Bytes(), []byte("Barangay2")) {
		t.Errorf("Expected barangay names in response, got %s", w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
