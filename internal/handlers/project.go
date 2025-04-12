package handlers

import (
	"net/http"
	"sync"
	"wow-bato-backend/internal/models"
	"wow-bato-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type ProjectHandlers struct {
	svc *services.ProjectService
}

func NewProjectHandlers(svc *services.ProjectService) *ProjectHandlers {
	return &ProjectHandlers{svc: svc}
}

func (h *ProjectHandlers) AddNewProject(c *gin.Context){

	session := services.CheckAuthentication(c)

	categoryID := c.Param("categoryID")
	barangayIDValue := session.Get("barangay_id")

	barangay_ID, ok := barangayIDValue.(uint)
	if !ok {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Invalid barangay_ID"})
		return
	}
	

	var newProject models.NewProject
	services.BindJSON(c, &newProject)

	err := h.svc.AddNewProject(barangay_ID, categoryID, newProject)
	services.CheckServiceError(c, err)

	c.IndentedJSON(http.StatusOK, gin.H{"message": "New Project Created"})
}

func (h *ProjectHandlers) DeleteProject(c *gin.Context){

	session := services.CheckAuthentication(c)

	projectID := c.Param("projectID")
	barangay_ID := session.Get("barangay_ID").(uint)

	err := h.svc.DeleteProject(barangay_ID, projectID)
	services.CheckServiceError(c, err)

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Project Deleted"})
}

func (h *ProjectHandlers) UpdateProject(c *gin.Context){

    session := services.CheckAuthentication(c)

    projectID := c.Param("projectID")
    barangay_ID := session.Get("barangay_ID").(uint)

    var updateProject models.UpdateProject
    services.BindJSON(c, &updateProject)

    err := h.svc.UpdateProject(barangay_ID, projectID, updateProject)
    services.CheckServiceError(c, err)

    c.IndentedJSON(http.StatusOK, gin.H{"message": "Updated Project"})
}

func (h *ProjectHandlers) GetAllProjects(c *gin.Context){

	session := services.CheckAuthentication(c)
	
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
		result, err := h.svc.GetAllProjects(barangay_ID, categoryID, limit, page)
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

func (h *ProjectHandlers) UpdateProjectStatus(c *gin.Context){

    session := services.CheckAuthentication(c)

    projectID := c.Param("projectID")
    barangay_ID := session.Get("barangay_ID").(uint)

    var newStatus models.NewProjectStatus
    services.BindJSON(c, &newStatus)

    err := h.svc.UpdateProjectStatus(projectID, barangay_ID, newStatus)
    services.CheckServiceError(c, err)
}

func (h *ProjectHandlers) GetSingleProject(c *gin.Context){

	services.CheckAuthentication(c)

	projectID := c.Param("projectID")

	project, err := h.svc.GetProjectSingle(projectID)
	services.CheckServiceError(c, err)

	c.IndentedJSON(http.StatusOK, gin.H{"data": project, "message": "Project " + project.Name +" Retrieved"})
}
