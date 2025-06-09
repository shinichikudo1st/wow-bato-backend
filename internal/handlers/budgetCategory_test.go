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
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestAddBudgetCategory(t *testing.T) {
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
	svc := services.NewBudgetCategoryService(gormDB)
	handlersObj := handlers.NewBudgetCategoryHandlers(svc)

	r.POST("/budget-category/add", func(c *gin.Context) {
		handlersObj.AddBudgetCategory(c)
	})

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "budget_categories"`).
		WithArgs(
			sqlmock.AnyArg(),          // gorm.Model fields (ID, CreatedAt, etc.)
			"Infra",                   // name
			"Infrastructure projects", // description
			1,                         // barangay_ID
		).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	newBudgetCategory := models.NewBudgetCategory{
		Name:        "Infra",
		Description: "Infrastructure projects",
		Barangay_ID: 1,
	}
	jsonValue, _ := json.Marshal(newBudgetCategory)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/budget-category/add", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("New Budget Category Added")) {
		t.Errorf("Expected success message, got %s", w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestDeleteBudgetCategory(t *testing.T) {
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
	svc := services.NewBudgetCategoryService(gormDB)
	handlersObj := handlers.NewBudgetCategoryHandlers(svc)

	r.DELETE("/budget-category/delete/:budget_ID", func(c *gin.Context) {
		handlersObj.DeleteBudgetCategory(c)
	})

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "budget_categories" WHERE "budget_categories"."id" = \$1`).
		WithArgs(2).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/budget-category/delete/2", nil)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("Budget Category Deleted")) {
		t.Errorf("Expected success message, got %s", w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestUpdateBudgetCategory(t *testing.T) {
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
	svc := services.NewBudgetCategoryService(gormDB)
	handlersObj := handlers.NewBudgetCategoryHandlers(svc)

	r.PUT("/budget-category/update/:budget_ID", func(c *gin.Context) {
		handlersObj.UpdateBudgetCategory(c)
	})

	// Mock the SELECT for existence check
	mock.ExpectQuery(`SELECT \* FROM "budget_categories" WHERE id = \$1 ORDER BY "budget_categories"."id" LIMIT 1`).
		WithArgs(2).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "barangay_ID"}).
			AddRow(2, "Old Name", "Old Description", 1))
	// Mock the UPDATE
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "budget_categories" SET (.+) WHERE "id" = \$1`).
		WithArgs("Updated Name", "Updated Description", 2).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	updateBudgetCategory := models.UpdateBudgetCategory{
		Name:        "Updated Name",
		Description: "Updated Description",
	}
	jsonValue, _ := json.Marshal(updateBudgetCategory)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/budget-category/update/2", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("Budget Category Updated")) {
		t.Errorf("Expected success message, got %s", w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestGetAllBudgetCategory(t *testing.T) {
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
	svc := services.NewBudgetCategoryService(gormDB)
	handlersObj := handlers.NewBudgetCategoryHandlers(svc)

	r.GET("/budget-category/all/:barangay_ID", func(c *gin.Context) {
		handlersObj.GetAllBudgetCategory(c)
	})

	// Mock the budget categories query
	mock.ExpectQuery(`SELECT budget_categories.id, budget_categories.name, budget_categories.description, budget_categories.barangay_ID, COUNT\(projects.id\) as project_count FROM "budget_categories" LEFT JOIN projects ON projects.category_id = budget_categories.id WHERE budget_categories.barangay_ID = \$1 GROUP BY budget_categories.id LIMIT \$2 OFFSET \$3`).
		WithArgs(1, 10, 0).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "barangay_ID", "project_count"}).
			AddRow(1, "Infra", "Infrastructure projects", 1, 2).
			AddRow(2, "Health", "Health projects", 1, 1))

	// Mock the count query
	mock.ExpectQuery(`SELECT count\(\*\) FROM "budget_categories" WHERE barangay_ID = \$1`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/budget-category/all/1?limit=10&page=1", nil)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("Infra")) || !bytes.Contains(w.Body.Bytes(), []byte("Health")) {
		t.Errorf("Expected budget category names in response, got %s", w.Body.String())
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("count")) {
		t.Errorf("Expected count in response, got %s", w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestGetSingleBudgetCategory(t *testing.T) {
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
	svc := services.NewBudgetCategoryService(gormDB)
	handlersObj := handlers.NewBudgetCategoryHandlers(svc)

	r.GET("/budget-category/single/:budget_ID", func(c *gin.Context) {
		// Simulate session with barangay_id
		session := sessions.Default(c)
		session.Set("barangay_id", uint(1))
		session.Save()
		handlersObj.GetSingleBudgetCategory(c)
	})

	mock.ExpectQuery(`SELECT \* FROM "budget_categories" WHERE barangay_ID = \$1 AND id = \$2`).
		WithArgs(1, 2).
		WillReturnRows(sqlmock.NewRows([]string{"name", "description"}).AddRow("Infra", "Infrastructure projects"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/budget-category/single/2", nil)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("Infra")) {
		t.Errorf("Expected budget category name in response, got %s", w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
