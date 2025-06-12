package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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

func TestAddNewProject(t *testing.T) {
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
	svc := services.NewProjectService(gormDB)
	handlersObj := handlers.NewProjectHandlers(svc, nil)

	r.POST("/project/add/:categoryID", func(c *gin.Context) {
		sess := sessions.Default(c)
		sess.Set("barangay_id", uint(1))
		sess.Save()
		handlersObj.AddNewProject(c)
	})

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "projects"`).
		WithArgs(
			1,                // barangay_id
			2,                // category_id
			"Test Project",   // name
			"A test project", // description
			"2024-06-01",     // start_date
			"2024-06-30",     // end_date
			"pending",        // status
		).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	newProject := models.NewProject{
		Name:        "Test Project",
		Description: "A test project",
		StartDate:   "2024-06-01",
		EndDate:     "2024-06-30",
		Status:      "pending",
	}
	jsonValue, _ := json.Marshal(newProject)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/project/add/2", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("New Project Created")) {
		t.Errorf("Expected success message, got %s", w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestDeleteProject(t *testing.T) {
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
	svc := services.NewProjectService(gormDB)
	handlersObj := handlers.NewProjectHandlers(svc, nil)

	r.DELETE("/project/delete/:projectID", func(c *gin.Context) {
		sess := sessions.Default(c)
		sess.Set("barangay_ID", uint(1))
		sess.Save()
		handlersObj.DeleteProject(c)
	})

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "projects" WHERE "projects"."id" = \$1`).
		WithArgs(2).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/project/delete/2", nil)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("Project Deleted")) {
		t.Errorf("Expected success message, got %s", w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestUpdateProject(t *testing.T) {
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
	svc := services.NewProjectService(gormDB)
	handlersObj := handlers.NewProjectHandlers(svc, nil)

	r.PUT("/project/update/:projectID", func(c *gin.Context) {
		sess := sessions.Default(c)
		sess.Set("barangay_ID", uint(1))
		sess.Save()
		handlersObj.UpdateProject(c)
	})

	mock.ExpectQuery(`SELECT \* FROM "projects" WHERE Barangay_ID = \$1 AND id = \$2`).
		WithArgs(1, 2).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "start_date", "end_date", "status", "barangay_id", "category_id"}).
			AddRow(2, "Old Name", "Old Description", "2024-06-01", "2024-06-30", "pending", 1, 1))
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "projects" SET (.+) WHERE "id" = \$1`).
		WithArgs("Updated Name", "Updated Description", 2).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	updateProject := models.UpdateProject{
		Name:        "Updated Name",
		Description: "Updated Description",
	}
	jsonValue, _ := json.Marshal(updateProject)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/project/update/2", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("Updated Project")) {
		t.Errorf("Expected success message, got %s", w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestGetAllProjects(t *testing.T) {
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
	svc := services.NewProjectService(gormDB)
	budgetSvc := services.NewBudgetCategoryService(gormDB)
	handlersObj := handlers.NewProjectHandlers(svc, budgetSvc)

	r.GET("/project/all/:categoryID", func(c *gin.Context) {
		sess := sessions.Default(c)
		sess.Set("barangay_id", uint(1))
		sess.Save()
		handlersObj.GetAllProjects(c)
	})

	mock.ExpectQuery(`SELECT id, name, status, start_date, end_date FROM "projects" WHERE barangay_id = \$1 AND category_id = \$2 LIMIT \$3 OFFSET \$4`).
		WithArgs(1, 2, 10, 0).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "status", "start_date", "end_date"}).
			AddRow(1, "Project 1", "pending", "2024-06-01", "2024-06-30").
			AddRow(2, "Project 2", "ongoing", "2024-07-01", "2024-07-31"))

	mock.ExpectQuery(`SELECT \* FROM "budget_categories" WHERE barangay_ID = \$1 AND id = \$2`).
		WithArgs(1, 2).
		WillReturnRows(sqlmock.NewRows([]string{"name", "description"}).AddRow("Infra", "Infrastructure projects"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/project/all/2?limit=10&page=1", nil)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("Project 1")) || !bytes.Contains(w.Body.Bytes(), []byte("Project 2")) {
		t.Errorf("Expected project names in response, got %s", w.Body.String())
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("Infra")) {
		t.Errorf("Expected budget category name in response, got %s", w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestUpdateProjectStatus(t *testing.T) {
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
	svc := services.NewProjectService(gormDB)
	handlersObj := handlers.NewProjectHandlers(svc, nil)

	r.PUT("/project/status/:projectID", func(c *gin.Context) {
		sess := sessions.Default(c)
		sess.Set("barangay_ID", uint(1))
		sess.Save()
		handlersObj.UpdateProjectStatus(c)
	})

	mock.ExpectQuery(`SELECT \* FROM "projects" WHERE Barangay_ID = \$1 AND id = \$2`).
		WithArgs(1, 2).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "start_date", "end_date", "status", "barangay_id", "category_id"}).
			AddRow(2, "Project Name", "Project Desc", "2024-06-01", "2024-06-30", "pending", 1, 1))
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "projects" SET (.+) WHERE "id" = \$1`).
		WithArgs("ongoing", "2024-07-01T00:00:00Z", 2).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	updateStatus := models.NewProjectStatus{
		Status:   "ongoing",
		FlexDate: time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC),
	}
	jsonValue, _ := json.Marshal(updateStatus)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/project/status/2", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
