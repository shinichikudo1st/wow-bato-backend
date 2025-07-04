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

func TestAddNewBudgetItem(t *testing.T) {
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
	svc := services.NewBudgetItemService(gormDB)
	handlersObj := handlers.NewBudgetItemHandlers(svc)

	r.POST("/budget-item/add/:projectID", func(c *gin.Context) {
		handlersObj.AddNewBudgetItem(c)
	})

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "budget_items"`).
		WithArgs(
			"Test Budget Item",   // name
			1000.0,               // amount_allocated
			"A test budget item", // description
			"pending",            // status
			sqlmock.AnyArg(),     // approval_date (nullable)
			1,                    // project_id
		).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	newBudgetItem := models.NewBudgetItem{
		Name:             "Test Budget Item",
		Amount_Allocated: 1000.0,
		Description:      "A test budget item",
		Status:           "pending",
	}
	jsonValue, _ := json.Marshal(newBudgetItem)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/budget-item/add/1", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("New Budget Item Added")) {
		t.Errorf("Expected success message, got %s", w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestGetAllBudgetItem(t *testing.T) {
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
	svc := services.NewBudgetItemService(gormDB)
	handlersObj := handlers.NewBudgetItemHandlers(svc)

	r.GET("/budget-item/all/:projectID", func(c *gin.Context) {
		handlersObj.GetAllBudgetItem(c)
	})

	mock.ExpectQuery(`SELECT \* FROM "budget_items" WHERE project_id = \$1 LIMIT \$2 OFFSET \$3`).
		WithArgs(1, 5, 0).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "amount_allocated", "description", "status", "approval_date", "project_id"}).
			AddRow(1, "Budget Item 1", 500.0, "Desc 1", "pending", nil, 1).
			AddRow(2, "Budget Item 2", 1000.0, "Desc 2", "approved", nil, 1))

	mock.ExpectQuery(`SELECT count\(\*\) FROM "budget_items" WHERE project_id = \$1`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/budget-item/all/1?filter=All&page=1", nil)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("Budget Item 1")) || !bytes.Contains(w.Body.Bytes(), []byte("Budget Item 2")) {
		t.Errorf("Expected budget item names in response, got %s", w.Body.String())
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("count")) {
		t.Errorf("Expected count in response, got %s", w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestGetSingleBudgetItem(t *testing.T) {
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
	svc := services.NewBudgetItemService(gormDB)
	handlersObj := handlers.NewBudgetItemHandlers(svc)

	r.GET("/budget-item/single/:projectID/:budgetItemID", func(c *gin.Context) {
		handlersObj.GetSingleBudgetItem(c)
	})

	mock.ExpectQuery(`SELECT \* FROM "budget_items" WHERE categoryID = \$1 AND status = \$2 ORDER BY "budget_items"."id" LIMIT 1`).
		WithArgs(1, 2).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "amount_allocated", "description", "status", "approval_date", "project_id"}).
			AddRow(2, "Budget Item 2", 1000.0, "Desc 2", "approved", nil, 1))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/budget-item/single/1/2", nil)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("Budget Item 2")) {
		t.Errorf("Expected budget item name in response, got %s", w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestUpdateStatusBudgetItem(t *testing.T) {
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
	svc := services.NewBudgetItemService(gormDB)
	handlersObj := handlers.NewBudgetItemHandlers(svc)

	r.PUT("/budget-item/status/:budgetItemID", func(c *gin.Context) {
		handlersObj.UpdateStatusBudgetItem(c)
	})

	mock.ExpectQuery(`SELECT \* FROM "budget_items" WHERE id = \$1 ORDER BY "budget_items"."id" LIMIT 1`).
		WithArgs(2).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "amount_allocated", "description", "status", "approval_date", "project_id"}).
			AddRow(2, "Budget Item 2", 1000.0, "Desc 2", "pending", nil, 1))
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "budget_items" SET (.+) WHERE "id" = \$1`).
		WithArgs("Approved", 2).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	updateStatus := models.UpdateStatus{
		Status: "approve",
	}
	jsonValue, _ := json.Marshal(updateStatus)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/budget-item/status/2", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("Budget Item Updated")) {
		t.Errorf("Expected success message, got %s", w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestDeleteBudgetItem(t *testing.T) {
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
	svc := services.NewBudgetItemService(gormDB)
	handlersObj := handlers.NewBudgetItemHandlers(svc)

	r.DELETE("/budget-item/delete/:budgetItemID", func(c *gin.Context) {
		handlersObj.DeleteBudgetItem(c)
	})

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "budget_items" WHERE "budget_items"."id" = \$1`).
		WithArgs(2).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/budget-item/delete/2", nil)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("Budget Item Deleted")) {
		t.Errorf("Expected success message, got %s", w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
