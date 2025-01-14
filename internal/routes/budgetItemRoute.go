package routes

import (
	"wow-bato-backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterBudgetItemRoutes registers budget item-related routes on the provided router group.
//
// This route contains:
// - POST /budgetItem/add/:categoryID: Handles adding a new budget item.
// - GET /budgetItem/all/:projectID: Handles retrieving a list of budget items.
// - GET /budgetItem/:projectID/:budgetItemID: Handles retrieving a single budget item.
// - PUT /budgetItem/:budgetItemID: Handles updating a budget item.
func RegisterBudgetItemRoutes(router *gin.RouterGroup) {
	budgetItem := router.Group("/budgetItem")
	{
		budgetItem.POST("/add/:categoryID", handlers.AddNewBudgetItem)
		budgetItem.GET("/all/:projectID", handlers.GetAllBudgetItem)
		budgetItem.GET("/:projectID/:budgetItemID", handlers.GetSingleBudgetItem)
		budgetItem.PUT("/:budgetItemID")
	}
}
