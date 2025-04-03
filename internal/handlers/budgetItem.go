package handlers

import (
	"net/http"
	"sync"
	"wow-bato-backend/internal/models"
	"wow-bato-backend/internal/services"

	"github.com/gin-gonic/gin"
)

func AddNewBudgetItem(c *gin.Context){
	
	services.CheckAuthentication(c)

	projectID := c.Param("projectID")

	var budgetItem models.NewBudgetItem
	services.BindJSON(c, &budgetItem)

	err := services.AddBudgetItem(projectID, budgetItem)
	services.CheckServiceError(c, err)

	c.IndentedJSON(http.StatusOK, gin.H{"message": "New Budget Item Added"})
}

func GetAllBudgetItem(c *gin.Context){
	
	services.CheckAuthentication(c)

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
	
	services.CheckAuthentication(c)

	projectID := c.Param("projectID")
	budgetItemID := c.Param("budgetItemID")

	budgetItem, err := services.GetSingleBudgetItem(projectID, budgetItemID)
	services.CheckServiceError(c, err)

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Retrieved Budget Items for category", "data": budgetItem})
}

func UpdateStatusBudgetItem(c *gin.Context){
	
	services.CheckAuthentication(c)

	budgetItemID := c.Param("budgetItemID")

	var newStatus models.UpdateStatus
	services.BindJSON(c, &newStatus)

	err := services.UpdateBudgetItemStatus(budgetItemID, newStatus)
	services.CheckServiceError(c, err)

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Budget Item Updated"})
}

func DeleteBudgetItem(c *gin.Context){
	
	services.CheckAuthentication(c)

	budgetItemID := c.Param("budgetItemID")

	err := services.DeleteBudgetItem(budgetItemID)
	services.CheckServiceError(c, err)

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Budget Item Deleted"})
}
