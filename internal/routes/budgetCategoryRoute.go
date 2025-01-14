package routes

import (
	"wow-bato-backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterBudgetCategoryRoutes registers budget category-related routes on the provided router group.
//
// This route contains:
// - POST /budgetCategory/add: Handles adding a new budget category.
// - DELETE /budgetCategory/delete/:budget_ID: Handles deleting a budget category.
// - PUT /budgetCategory/update/:budget_ID: Handles updating a budget category.
// - GET /budgetCategory/all/:barangay_ID: Handles retrieving a list of budget categories.
// - GET /budgetCategory/:barangay_ID/:budget_ID: Handles retrieving a single budget category.
func RegisterBudgetCategoryRoutes(router *gin.RouterGroup) {
	budgetCategory := router.Group("/budgetCategory")
	{
		budgetCategory.POST("/add", handlers.AddBudgetCategory)
		budgetCategory.DELETE("/delete/:budget_ID", handlers.DeleteBudgetCategory)
		budgetCategory.PUT("/update/:budget_ID", handlers.UpdateBudgetCategory)
		budgetCategory.GET("/all/:barangay_ID", handlers.GetAllBudgetCategory)
		budgetCategory.GET("/:barangay_ID/:budget_ID", handlers.GetSingleBudgetCategory)
	}
}
