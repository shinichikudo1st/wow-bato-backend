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
	status := c.Param("status")

	budgetItem, err := services.GetAllBudgetItem(categoryID, status)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Retrieved Budget Items for category", "data": budgetItem})
}
