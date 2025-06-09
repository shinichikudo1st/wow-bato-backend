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
