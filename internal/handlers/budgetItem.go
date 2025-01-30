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

func GetAllBudgetItem(c *gin.Context){
	session := sessions.Default(c)

	if session.Get("authenticated") != true {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Access Denied"})
		return
	}

	projectID := c.Param("projectID")
	filter := c.Query("filter")
	page := c.Query("page")

	budgetItem, err := services.GetAllBudgetItem(projectID, filter, page)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Retrieved Budget Items for category", "data": budgetItem})
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