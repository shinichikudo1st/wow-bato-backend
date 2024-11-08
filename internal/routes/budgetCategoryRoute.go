package routes

import (
	"wow-bato-backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterBudgetCategoryRoutes(router *gin.RouterGroup){
	budgetCategory := router.Group("/budgetCategory")
	{
		budgetCategory.POST("/add", handlers.AddBudgetCategory)
		budgetCategory.DELETE("/delete/:budget_ID", handlers.DeleteBudgetCategory)
		budgetCategory.PUT("/update/:budget_ID", handlers.UpdateBudgetCategory)
		budgetCategory.GET("/all/:barangay_ID", handlers.GetAllBudgetCategory)
		budgetCategory.GET("/:barangay_ID/:budget_ID", handlers.GetSingleBudgetCategory)
	}
}