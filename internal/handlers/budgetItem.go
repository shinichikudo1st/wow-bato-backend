package handlers

import (
	"net/http"
	"wow-bato-backend/internal/models"
	"wow-bato-backend/internal/services"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AddNewBudgetItem(c *gin.Context){
	session := sessions.Default(c)

	if session.Get("authenticated") != true {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Access Denied"})
		return
	}

	categoryID := c.Param("categoryID")

	var budgetItem models.NewBudgetItem
	if err := c.ShouldBindJSON(&budgetItem); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := services.AddBudgetItem(categoryID, budgetItem)
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

	categoryID := c.Param("categoryID")
	filter := c.Query("filter")

	budgetItem, err := services.GetAllBudgetItem(categoryID, filter)

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

	categoryID := c.Param("categoryID")
	budgetItemID := c.Param("budgetItemID")

	budgetItem, err := services.GetSingleBudgetItem(categoryID, budgetItemID)

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