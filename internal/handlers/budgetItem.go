package handlers

import (
	"net/http"
	"sync"
	"wow-bato-backend/internal/models"
	"wow-bato-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type BudgetItemHandlers struct {
	svc *services.BudgetItemService
}

func NewBudgetItemHandlers(svc *services.BudgetItemService) *BudgetItemHandlers {
	return &BudgetItemHandlers{svc: svc}
}

func (h *BudgetItemHandlers) AddNewBudgetItem(c *gin.Context){
	
	services.CheckAuthentication(c)

	projectID := c.Param("projectID")

	var budgetItem models.NewBudgetItem
	services.BindJSON(c, &budgetItem)

	err := h.svc.AddBudgetItem(projectID, budgetItem)
	services.CheckServiceError(c, err)

	c.IndentedJSON(http.StatusOK, gin.H{"message": "New Budget Item Added"})
}

func (h *BudgetItemHandlers) GetAllBudgetItem(c *gin.Context){
	
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
		result, err := h.svc.GetAllBudgetItem(projectID, filter, page)
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
		result, err := h.svc.CountBudgetItem(projectID)
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

func (h *BudgetItemHandlers) GetSingleBudgetItem(c *gin.Context){
	
	services.CheckAuthentication(c)

	projectID := c.Param("projectID")
	budgetItemID := c.Param("budgetItemID")

	budgetItem, err := h.svc.GetSingleBudgetItem(projectID, budgetItemID)
	services.CheckServiceError(c, err)

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Retrieved Budget Items for category", "data": budgetItem})
}

func (h *BudgetItemHandlers) UpdateStatusBudgetItem(c *gin.Context){
	
	services.CheckAuthentication(c)

	budgetItemID := c.Param("budgetItemID")

	var newStatus models.UpdateStatus
	services.BindJSON(c, &newStatus)

	err := h.svc.UpdateBudgetItemStatus(budgetItemID, newStatus)
	services.CheckServiceError(c, err)

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Budget Item Updated"})
}

func (h *BudgetItemHandlers) DeleteBudgetItem(c *gin.Context){
	
	services.CheckAuthentication(c)

	budgetItemID := c.Param("budgetItemID")

	err := h.svc.DeleteBudgetItem(budgetItemID)
	services.CheckServiceError(c, err)

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Budget Item Deleted"})
}
