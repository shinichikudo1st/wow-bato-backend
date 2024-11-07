package routes

import (
	"wow-bato-backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterBudgetCategoryRoutes(router *gin.RouterGroup){
	budgetCategory := router.Group("/budgetCategory")
	{
		budgetCategory.POST("/add", handlers.AddBudgetCategory)
		budgetCategory.DELETE("/delete/:budgetID", handlers.DeleteBudgetCategory)
	}
}