package routes

import (
	"wow-bato-backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterBudgetItemRoutes(router *gin.RouterGroup) {
	budgetItem := router.Group("/budgetItem")
	{
		budgetItem.POST("/add/:projectID", handlers.AddNewBudgetItem)
		budgetItem.GET("/all/:projectID", handlers.GetAllBudgetItem)
		budgetItem.GET("/:projectID/:budgetItemID", handlers.GetSingleBudgetItem)
		budgetItem.PUT("/update-status/:budgetItemID", handlers.UpdateStatusBudgetItem)
		budgetItem.DELETE("/delete/:budgetItemID", handlers.DeleteBudgetItem)
	}
}
