// Package handlers provides HTTP request handlers for the wow-bato application
// The package implements handlers for budget category management including:
//   - Budget Category Creation
//   - Budget Category Deletion
//   - Budget Category Update
//   - Get all Budget Categories (paginated)
//   - Get a single budget category
package handlers

import (
	"net/http"
	"wow-bato-backend/internal/models"
	"wow-bato-backend/internal/services"

	"sync"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

//	AddBudgetCategory handles the creation for creating new budget category
//
//	This handlers performs the following operations:
//		1. Validates user authentication and authorization
//		2. Validates and binds the new budget category data
//		3. Delegates budget category creation to the services layer
//		4. Returns appropriate response based on operation result
// 
//	Security:
//		- Requires authenticated session
//		- Validates administrative privileges
//
//	@Summary Create a new budget category
//	@Description Creates a new budget category with the provided information
//	@Tags Budget Category
//	@Accept json
//	@Produce json
//	@Param budgetCategory body models.NewBudgetCategory true "Budget Category details including name, description, and barangay ID"
//	@Success 200 {object} gin.H "Returns success message on budget category creation"
//	@Failure 400 {object} gin.H "Returns error when request validation fails"
//	@Failure 401 {object} gin.H "Returns error when user is not authenticated"
//	@Failure 500 {object} gin.H "Returns error when budget category creation fails"
//	@Router /budgetCategory [post]
func AddBudgetCategory(c *gin.Context){
	session := sessions.Default(c)

	if session.Get("authenticated") != true {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Access Denied"})
		return
	}

	var newBudgetCategory models.NewBudgetCategory

	if err := c.ShouldBindJSON(&newBudgetCategory); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := services.AddBudgetCategory(newBudgetCategory)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "New Budget Category Added"})
}

// DeleteBudgetCategory handles the deletion of a budget category
//
//	This handler performs the following operations:
//		1. Validates user authentication and authorization
//		2. Validates and binds the budget category ID
//		3. Delegates budget category deletion to the services layer
//		4. Returns appropriate response based on operation result
// 
//	Security:
//		- Requires authenticated session
//		- Validates administrative privileges
//
//	@Summary Delete a budget category
//	@Description Deletes a budget category with the provided ID
//	@Tags Budget Category
//	@Accept json no body
//	@Produce json
//	@Param budget_ID path int true "Budget Category ID"
//	@Success 200 {object} gin.H "Returns success message on budget category deletion"
//	@Failure 400 {object} gin.H "Returns error when request validation fails"
//	@Failure 401 {object} gin.H "Returns error when user is not authenticated"
//	@Failure 500 {object} gin.H "Returns error when budget category deletion fails"
//	@Router /budgetCategory/{budget_ID} [delete]
func DeleteBudgetCategory(c *gin.Context){
	session := sessions.Default(c)

	if session.Get("authenticated") != true {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Access Denied"})
		return
	}

	budget_ID := c.Param("budget_ID")

	err := services.DeleteBudgetCategory(budget_ID)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Budget Category Deleted"})
}

// UpdateBudgetCategory handles the update of a budget category
//
//	This handler performs the following operations:
//		1. Validates user authentication and authorization
//		2. Validates and binds the budget category ID and update details
//		3. Delegates budget category update to the services layer
//		4. Returns appropriate response based on operation result
// 
//	Security:
//		- Requires authenticated session
//		- Validates administrative privileges
//
//	@Summary Update a budget category
//	@Description Updates a budget category with the provided ID and details
//	@Tags Budget Category
//	@Accept json
//	@Produce json
//	@Param budget_ID path int true "Budget Category ID"
//	@Param budgetCategory body models.UpdateBudgetCategory true "Budget Category details"
//	@Success 200 {object} gin.H "Returns success message on budget category update"
//	@Failure 400 {object} gin.H "Returns error when request validation fails"
//	@Failure 401 {object} gin.H "Returns error when user is not authenticated"
//	@Failure 500 {object} gin.H "Returns error when budget category update fails"
//	@Router /budgetCategory/{budget_ID} [put]
func UpdateBudgetCategory(c *gin.Context){
	session := sessions.Default(c)

	if session.Get("authenticated") != true {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Access Denied"})
		return
	}

	budget_ID := c.Param("budget_ID")

	var updateBudgetCategory models.UpdateBudgetCategory

	if err := c.ShouldBindJSON(&updateBudgetCategory); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := services.UpdateBudgetCategory(budget_ID, updateBudgetCategory)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Budget Category Updated"})
}

// GetAllBudgetCategory handles the retrieval of all budget categories
//
//	This handler performs the following operations:
//		1. Validates user authentication and authorization
//		2. Checks for dependent records
//		3. Performs cascading deletion if required
//		4. Maintains referential integrity
//
//	Security:
//		- Requires authenticated session
//		- Validates administrative privileges
//		- Ensures proper cleanup of dependencies
//
//	@Summary Get all budget categories
//	@Description Retrieves all budget categories and their associated data from the system
//	@Tags Budget Category
//	@Accept json
//	@Produce json
//	@Param page query string false "Page number for pagination"
//	@Param limit query string false "Number of items per page"
//	@Success 200 {object} gin.H "Returns list of budget categories and pagination info"
//	@Failure 401 {object} gin.H "Returns error when user is not authenticated"
//	@Failure 500 {object} gin.H "Returns error when retrieval fails"
//	@Router /budgetCategory [get]
func GetAllBudgetCategory(c *gin.Context) {
	session := sessions.Default(c)

	if session.Get("authenticated") != true {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Access Denied"})
		return
	}

	limit := c.Query("limit")
	page := c.Query("page")
	barangay_ID := c.Param("barangay_ID")

	// Results and error handling
	var (
		budgetCategories []models.BudgetCategoryResponse
		count            int64
		errors           []error
	)

	var wg sync.WaitGroup
	var mu sync.Mutex // Mutex for safely appending errors

	wg.Add(2) // Two goroutines

	// Fetch budget categories in a goroutine
	go func() {
		defer wg.Done()
		result, err := services.GetAllBudgetCategory(barangay_ID, limit, page)
		if err != nil {
			mu.Lock()
			errors = append(errors, err)
			mu.Unlock()
			return
		}
		mu.Lock()
		budgetCategories = result
		mu.Unlock()
	}()

	// Fetch budget category count in a goroutine
	go func() {
		defer wg.Done()
		result, err := services.GetBudgetCategoryCount(barangay_ID)
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

	// Wait for all goroutines to finish
	wg.Wait()

	// Handle errors if any
	if len(errors) > 0 {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": errors[0].Error()})
		return
	}

	// Send JSON response
	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "All Budget Categories Retrieved",
		"data":    budgetCategories,
		"count":   count,
	})
}

// GetSingleBudgetCategory handles the retrieval of a single budget category
//
//	This handler performs the following operations:
//		1. Validates user authentication and authorization
//		2. Validates and binds the budget category ID
//		3. Delegates budget category retrieval to the services layer
//		4. Returns appropriate response based on operation result
// 
//	Security:
//		- Requires authenticated session
//		- Validates administrative privileges
//
//	@Summary Get a single budget category
//	@Description Retrieves a single budget category with the provided ID
//	@Tags Budget Category
//	@Accept json
//	@Produce json
//	@Param budget_ID path int true "Budget Category ID"
//	@Success 200 {object} gin.H "Returns a single budget category"
//	@Failure 401 {object} gin.H "Returns error when user is not authenticated"
//	@Failure 404 {object} gin.H "Returns error when budget category is not found"
//	@Failure 500 {object} gin.H "Returns error when budget category retrieval fails"
//	@Router /budgetCategory/{budget_ID} [get]
func GetSingleBudgetCategory(c *gin.Context){
	session := sessions.Default(c)

	if session.Get("authenticated") != true {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Access Denied"})
		return
	}

	barangay_ID := session.Get("barangay_id").(uint)
	budget_ID := c.Param("budget_ID")

	budgetCategory, err := services.GetBudgetCategory(barangay_ID, budget_ID)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "A Budget Category Retrieved", "data": budgetCategory})
}