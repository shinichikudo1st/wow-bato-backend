// Package handlers provides HTTP request handlers for the wow-bato application.
// It implements handlers for project management operations, including:
//   - Project creation and initialization
//   - Project updates and modifications
//   - Project status management
//   - Project listing and retrieval
//   - Project deletion and cleanup
//
// The package ensures proper authentication and authorization checks
// while maintaining data consistency across project operations.
package handlers

import (
	"net/http"
	"sync"
	"wow-bato-backend/internal/models"
	"wow-bato-backend/internal/services"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// AddNewProject handles the creation of new projects within a budget category.
//
// This handler performs the following operations:
//  1. Validates user authentication and session
//  2. Extracts barangay ID from session and category ID from URL
//  3. Validates and binds the project data
//  4. Delegates project creation to the services layer
//
// Security:
//  - Requires authenticated session
//  - Validates barangay association
//  - Enforces proper authorization
//
// @Summary Create a new project in a budget category
// @Description Creates a new project with the provided details and associates it with a budget category
// @Tags Project
// @Accept json
// @Produce json
// @Param categoryID path string true "ID of the budget category"
// @Param project body models.NewProject true "Project details including name, description, and budget"
// @Success 200 {object} gin.H "Returns success message on project creation"
// @Failure 400 {object} gin.H "Returns error when request validation fails"
// @Failure 401 {object} gin.H "Returns error when user is not authenticated"
// @Failure 500 {object} gin.H "Returns error when project creation fails"
// @Router /projects/{categoryID} [post]
func AddNewProject(c *gin.Context){

	session := sessions.Default(c)
	services.CheckAuthentication(c, session)

	categoryID := c.Param("categoryID")
	barangayIDValue := session.Get("barangay_id")

	barangay_ID, ok := barangayIDValue.(uint)
	if !ok {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Invalid barangay_ID"})
		return
	}
	

	var newProject models.NewProject
	services.BindJSON(c, &newProject)

	err := services.AddNewProject(barangay_ID, categoryID, newProject)
	services.CheckServiceError(c, err)

	c.IndentedJSON(http.StatusOK, gin.H{"message": "New Project Created"})
}

// DeleteProject handles the removal of existing projects.
//
// This handler performs the following operations:
//  1. Validates user authentication and authorization
//  2. Extracts project ID from the request
//  3. Verifies project ownership and access rights
//  4. Executes project deletion with proper cleanup
//
// Security:
//  - Requires authenticated session
//  - Validates project ownership
//  - Ensures proper access rights
//
// @Summary Delete an existing project
// @Description Removes a project and its associated data from the system
// @Tags Project
// @Accept json
// @Produce json
// @Param projectID path string true "ID of the project to delete"
// @Success 200 {object} gin.H "Returns success message on project deletion"
// @Failure 401 {object} gin.H "Returns error when user is not authenticated"
// @Failure 403 {object} gin.H "Returns error when user lacks permission"
// @Failure 500 {object} gin.H "Returns error when deletion fails"
// @Router /projects/{projectID} [delete]
func DeleteProject(c *gin.Context){

	session := sessions.Default(c)
	services.CheckAuthentication(c, session)

	projectID := c.Param("projectID")
	barangay_ID := session.Get("barangay_ID").(uint)

	err := services.DeleteProject(barangay_ID, projectID)
	services.CheckServiceError(c, err)

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Project Deleted"})
}

// UpdateProject handles modifications to existing project details.
//
// This handler performs the following operations:
//  1. Validates user authentication and authorization
//  2. Extracts project ID and update data
//  3. Validates update permissions
//  4. Applies and persists project changes
//
// The update operation supports modifying:
//  - Project name and description
//  - Budget allocation
//  - Project metadata
//
// @Summary Update project information
// @Description Modifies existing project details with provided updates
// @Tags Project
// @Accept json
// @Produce json
// @Param projectID path string true "ID of the project to update"
// @Param project body models.UpdateProject true "Updated project details"
// @Success 200 {object} gin.H "Returns success message on update"
// @Failure 400 {object} gin.H "Returns error when request validation fails"
// @Failure 401 {object} gin.H "Returns error when user is not authenticated"
// @Failure 500 {object} gin.H "Returns error when update fails"
// @Router /projects/{projectID} [put]
func UpdateProject(c *gin.Context){

    session := sessions.Default(c)
    services.CheckAuthentication(c, session)

    projectID := c.Param("projectID")
    barangay_ID := session.Get("barangay_ID").(uint)

    var updateProject models.UpdateProject
    services.BindJSON(c, &updateProject)

    err := services.UpdateProject(barangay_ID, projectID, updateProject)
    services.CheckServiceError(c, err)

    c.IndentedJSON(http.StatusOK, gin.H{"message": "Updated Project"})
}

// GetAllProjects retrieves a paginated list of projects with optional filtering.
//
// This handler performs the following operations:
//  1. Validates user authentication
//  2. Processes pagination parameters
//  3. Applies any specified filters
//  4. Retrieves and formats project data
//
// The response includes:
//  - Project basic information
//  - Associated budget category
//  - Current status and progress
//  - Creation and modification timestamps
//
// @Summary Retrieve all projects with pagination
// @Description Gets a list of projects with optional filtering and pagination
// @Tags Project
// @Accept json
// @Produce json
// @Param page query string false "Page number for pagination"
// @Param limit query string false "Number of items per page"
// @Success 200 {object} gin.H "Returns list of projects and pagination info"
// @Failure 401 {object} gin.H "Returns error when user is not authenticated"
// @Failure 500 {object} gin.H "Returns error when retrieval fails"
// @Router /projects [get]
func GetAllProjects(c *gin.Context){

	session := sessions.Default(c)
	services.CheckAuthentication(c, session)
	
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

// UpdateProjectStatus handles project status transitions and updates.
//
// This handler performs the following operations:
//  1. Validates user authentication and authorization
//  2. Verifies valid status transition
//  3. Updates project status with timestamp
//  4. Records status change history
//
// Status updates include:
//  - New status designation
//  - Timestamp of change
//  - Optional status notes or comments
//
// @Summary Update project status
// @Description Changes the current status of a project
// @Tags Project
// @Accept json
// @Produce json
// @Param projectID path string true "ID of the project"
// @Param status body models.NewProjectStatus true "New status details"
// @Success 200 {object} gin.H "Returns success message on status update"
// @Failure 400 {object} gin.H "Returns error when request validation fails"
// @Failure 401 {object} gin.H "Returns error when user is not authenticated"
// @Failure 500 {object} gin.H "Returns error when status update fails"
// @Router /projects/{projectID}/status [put]
func UpdateProjectStatus(c *gin.Context){

    session := sessions.Default(c)
    services.CheckAuthentication(c, session)

    projectID := c.Param("projectID")
    barangay_ID := session.Get("barangay_ID").(uint)

    var newStatus models.NewProjectStatus
    services.BindJSON(c, &newStatus)

    err := services.UpdateProjectStatus(projectID, barangay_ID, newStatus)
    services.CheckServiceError(c, err)
}

func GetSingleProject(c *gin.Context){

	session := sessions.Default(c)
	services.CheckAuthentication(c, session)

	projectID := c.Param("projectID")

	project, err := services.GetProjectSingle(projectID)
	services.CheckServiceError(c, err)

	c.IndentedJSON(http.StatusOK, gin.H{"data": project, "message": "Project " + project.Name +" Retrieved"})
}
