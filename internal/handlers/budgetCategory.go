package handlers

import (
	"net/http"
	"wow-bato-backend/internal/models"
	"wow-bato-backend/internal/services"

	"sync"

	"github.com/gin-gonic/gin"
)

func AddBudgetCategory(c *gin.Context){
	
	services.CheckAuthentication(c)

	var newBudgetCategory models.NewBudgetCategory
	services.BindJSON(c, &newBudgetCategory)

	err := services.AddBudgetCategory(newBudgetCategory)
	services.CheckServiceError(c, err)

	c.IndentedJSON(http.StatusOK, gin.H{"message": "New Budget Category Added"})
}

func DeleteBudgetCategory(c *gin.Context){

	services.CheckAuthentication(c)

	budget_ID := c.Param("budget_ID")

	err := services.DeleteBudgetCategory(budget_ID)
	services.CheckServiceError(c, err)

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Budget Category Deleted"})
}

func UpdateBudgetCategory(c *gin.Context){

	services.CheckAuthentication(c)

	budget_ID := c.Param("budget_ID")

	var updateBudgetCategory models.UpdateBudgetCategory
	services.BindJSON(c, &updateBudgetCategory)

	err := services.UpdateBudgetCategory(budget_ID, updateBudgetCategory)
	services.CheckServiceError(c, err)

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Budget Category Updated"})
}

func GetAllBudgetCategory(c *gin.Context) {
	
	services.CheckAuthentication(c)

	limit := c.Query("limit")
	page := c.Query("page")
	barangay_ID := c.Param("barangay_ID")

	var (
		budgetCategories []models.BudgetCategoryResponse
		count            int64
		errors           []error
	)

	var wg sync.WaitGroup
	var mu sync.Mutex 

	wg.Add(2) 

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

	wg.Wait()

	if len(errors) > 0 {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": errors[0].Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "All Budget Categories Retrieved",
		"data":    budgetCategories,
		"count":   count,
	})
}

func GetSingleBudgetCategory(c *gin.Context){
	
	session := services.CheckAuthentication(c)

	barangay_ID := session.Get("barangay_id").(uint)
	budget_ID := c.Param("budget_ID")

	budgetCategory, err := services.GetBudgetCategory(barangay_ID, budget_ID)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "A Budget Category Retrieved", "data": budgetCategory})
}