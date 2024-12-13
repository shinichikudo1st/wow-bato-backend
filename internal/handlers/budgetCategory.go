package handlers

import (
	"net/http"
	"wow-bato-backend/internal/models"
	"wow-bato-backend/internal/services"

	"sync"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

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

func GetSingleBudgetCategory(c *gin.Context){
	session := sessions.Default(c)

	if session.Get("authenticated") != true {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Access Denied"})
		return
	}

	barangay_ID := c.Param("barangay_ID")
	budget_ID := c.Param("budget_ID")

	budgetCategory, err := services.GetBudgetCategory(barangay_ID, budget_ID)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "A Budget Category Retrieved", "data": budgetCategory})
}