package routes

import (
	"wow-bato-backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterBudgetCategoryRoutes(router *gin.RouterGroup){
	budgetCategory := router.Group("/budgetCategory")
	{
		budgetCategory.POST("/addBudgetCategory", handlers.AddBudgetCategory)
		budgetCategory.DELETE("/deleteBudgetCategory/:budgetID", handlers.DeleteBudgetCategory)
	}
}