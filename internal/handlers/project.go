package handlers

import (
	"net/http"
	"sync"
	"wow-bato-backend/internal/models"
	"wow-bato-backend/internal/services"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Handler for adding new project
// @Summary Add new project
// @Tags Project
// @Accept json
// @Produce json
// @Param project body models.NewProject true "Project details"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
func AddNewProject(c *gin.Context){
	session := sessions.Default(c)

	if session.Get("authenticated") != true {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Access Denied"})
		return
	}

	categoryID := c.Param("categoryID")
	barangayIDValue := session.Get("barangay_id")
	barangay_ID, ok := barangayIDValue.(uint)
	if !ok {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Invalid barangay_ID"})
		return
	}
	

	var newProject models.NewProject
	if err := c.ShouldBindJSON(&newProject); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := services.AddNewProject(barangay_ID, categoryID, newProject)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "New Project Created"})
}

// Handler for deleting project
// @Summary Delete project
// @Tags Project
// @Accept json no body
// @Produce json
// @Success 200 {object} gin.H
// @Failure 500 {object} gin.H
func DeleteProject(c *gin.Context){
	session := sessions.Default(c)

	if session.Get("authenticated") != true {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Access Denied"})
		return
	}

	projectID := c.Param("projectID")
	barangay_ID := session.Get("barangay_ID").(uint)

	err := services.DeleteProject(barangay_ID, projectID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Project Deleted"})
}

// Handler for updating project
// @Summary Update project
// @Tags Project
// @Accept json
// @Produce json
// @Param project body models.UpdateProject true "Project details"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
func UpdateProject(c *gin.Context){
    session := sessions.Default(c)
    if session.Get("authenticated") != true {
        c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Access Denied"})
        return
    }

    projectID := c.Param("projectID")
    barangay_ID := session.Get("barangay_ID").(uint)

    var updateProject models.UpdateProject
    if err := c.ShouldBindJSON(&updateProject); err != nil {
        c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    err := services.UpdateProject(barangay_ID, projectID, updateProject)
    if err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.IndentedJSON(http.StatusOK, gin.H{"message": "Updated Project"})
}

// Handler for getting all projects
// @Summary Get all projects
// @Tags Project
// @Accept json no body
// @Produce json
// @Success 200 {object} gin.H
// @Failure 500 {object} gin.H
func GetAllProjects(c *gin.Context){
	session := sessions.Default(c)

	if session.Get("authenticated") != true {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Access Denied"})
		return
	}
	
	categoryID := c.Param("categoryID")
	barangay_ID, ok := session.Get("barangay_id").(uint)
	if !ok {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Invalid session data"})
		return
	}
	limit := c.Query("limit")
	page := c.Query("page")

	var (
		projectList []models.ProjectList
		budgetCategory	models.DisplayBudgetCategory
		errors		[]error
	)

	var wg sync.WaitGroup
	var mu sync.Mutex

	wg.Add(2)

	go func(){
		defer wg.Done()
		result, err := services.GetAllProjects(barangay_ID, categoryID, limit, page)
		if err != nil {
			mu.Lock()
			errors = append(errors, err)
			mu.Unlock()
			return
		}
		mu.Lock()
		projectList = result
		mu.Unlock()
	}()

	go func(){
		defer wg.Done()
		result, err := services.GetBudgetCategory(barangay_ID, categoryID)
		if err != nil {
			mu.Lock()
			errors = append(errors, err)
			mu.Unlock()
			return
		}
		mu.Lock()
		budgetCategory = result
		mu.Unlock()
	}()
	
	wg.Wait()

	if len(errors) > 0 {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": errors[0].Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"projects": projectList, "category": budgetCategory})
}

// Handler for updating project status
// @Summary Update project status
// @Tags Project
// @Accept json
// @Produce json
// @Param projectStatus body models.NewProjectStatus true "Project status details"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
func UpdateProjectStatus(c *gin.Context){
    session := sessions.Default(c)
    if session.Get("authenticated") != true {
        c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Access Denied"})
        return 
    }

    projectID := c.Param("projectID")
    barangay_ID := session.Get("barangay_ID").(uint)

    var newStatus models.NewProjectStatus
    if err := c.ShouldBindJSON(&newStatus); err != nil {
        c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    err := services.UpdateProjectStatus(projectID, barangay_ID, newStatus)
    if err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
}
