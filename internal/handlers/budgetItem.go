// Package handlers provides HTTP request handlers for the wow-bato application
// The package implements handlers for budget category management including:
//   - Budget Item Creation
//   - Budget Item Deletion
//   - Budget Item Update
//   - Get all Budget Items (paginated)
//   - Get a single budget item
package handlers

import (
	"net/http"
	"sync"
	"wow-bato-backend/internal/models"
	"wow-bato-backend/internal/services"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

//	AddNewBudgetItem handles the creation for creating new budget item
//
//	This handlers performs the following operations:
//		1. Validates user authentication and authorization
//		2. Validates and binds the new budget item data
//		3. Delegates budget item creation to the services layer
//		4. Returns appropriate response based on operation result
//
//	Security:
//		- Requires authenticated session
//		- Validates administrative privileges
//
//	@Summary Create a new budget item
//	@Description Creates a new budget item with the provided information
//	@Tags Budget Item
//	@Accept json
//	@Produce json
//	@Param budgetItem body models.NewBudgetItem true "Budget Item details including name, description, and barangay ID"
//	@Success 200 {object} gin.H "Returns success message on budget item creation"
//	@Failure 400 {object} gin.H "Returns error when request validation fails"
//	@Failure 401 {object} gin.H "Returns error when user is not authenticated"
//	@Failure 500 {object} gin.H "Returns error when budget item creation fails"
//	@Router /budgetItem [post]
func AddNewBudgetItem(c *gin.Context){
	session := sessions.Default(c)

	if session.Get("authenticated") != true {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Access Denied"})
		return
	}

	projectID := c.Param("projectID")

	var budgetItem models.NewBudgetItem
	if err := c.ShouldBindJSON(&budgetItem); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := services.AddBudgetItem(projectID, budgetItem)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "New Budget Item Added"})
}

//	GetAllBudgetItem handles the retrieval for budget items and count based on the filter
//
//	This handlers performs the following operations:
//		1. Validates user authentication and authorization
//		2. Collects the necessary URL parameter(projectID) and Query parameters(filter and page)
//		3. Delegates budget item retrieval to the services layer
//		4. Delegates counting of budget items to the services layer
//		5. Returns budget item slices response as JSON data to client
//
//	Security:
//		- Requires authenticated session
//		- Validates administrative privileges
//
//	@Summary Retrieves budget item
//	@Description Retrieves budget items and count based on the page and filter
//	@Tags Budget Item
//	@No accepted json
//	@Produce json
//	@Param Query param: page and filter, Path param: projectID
//	@Success 200 {object} gin.H "Returns success message and budget item slices as json on budget item retrieval"
//	@Failure 400 {object} gin.H "Returns error when request validation fails"
//	@Failure 401 {object} gin.H "Returns error when user is not authenticated"
//	@Failure 500 {object} gin.H "Returns error when budget retrieval fails"
//	@Router /budgetItem [get]
func GetAllBudgetItem(c *gin.Context){
	session := sessions.Default(c)

	if session.Get("authenticated") != true {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Access Denied"})
		return
	}

	projectID := c.Param("projectID")
	filter := c.Query("filter")
	page := c.Query("page")

	var (
		budgetItems []models.Budget_Item
		count		int64
		errors		[]error
	)

	var wg sync.WaitGroup
	var mu sync.Mutex

	wg.Add(2)

	go func ()  {
		defer wg.Done()
		result, err := services.GetAllBudgetItem(projectID, filter, page)
		if err != nil {
			mu.Lock()
			errors = append(errors, err)
			mu.Unlock()
			return
		}
		mu.Lock()
		budgetItems = result
		mu.Unlock()
	}()

	go func () {
		defer wg.Done()
		result, err := services.CountBudgetItem(projectID)
		if err != nil {
			mu.Lock()
			errors = append(errors, err)
			mu.Unlock()
			return
		}
		mu.Lock()
		count = result
		mu.Unlock()
	}()

	wg.Wait()

	if len(errors) > 0 {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": errors[0].Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Retrieved Budget Items for category", "data": budgetItems, "count": count})
}

func GetSingleBudgetItem(c *gin.Context){
	session := sessions.Default(c)

	if session.Get("authenticated") != true {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Access Denied"})
		return
	}

	projectID := c.Param("projectID")
	budgetItemID := c.Param("budgetItemID")

	budgetItem, err := services.GetSingleBudgetItem(projectID, budgetItemID)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Retrieved Budget Items for category", "data": budgetItem})
}

func UpdateStatusBudgetItem(c *gin.Context){
	session := sessions.Default(c)

	if session.Get("authenticated") != true {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Access Denied"})
		return
	}

	budgetItemID := c.Param("budgetItemID")

	var newStatus models.UpdateStatus
	if err := c.ShouldBindJSON(&newStatus); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	err := services.UpdateBudgetItemStatus(budgetItemID, newStatus)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Budget Item"})
}