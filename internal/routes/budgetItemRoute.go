package routes

import (
	"wow-bato-backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterBudgetItemRoutes(router *gin.RouterGroup){
	budgetItem := router.Group("/budgetItem")
	{
		budgetItem.POST("/add/:categoryID", handlers.AddNewBudgetItem)
		budgetItem.GET("/all/:projectID", handlers.GetAllBudgetItem)
		budgetItem.GET("/:projectID/:budgetItemID", handlers.GetSingleBudgetItem)
		budgetItem.PUT("/:budgetItemID")
	}
}