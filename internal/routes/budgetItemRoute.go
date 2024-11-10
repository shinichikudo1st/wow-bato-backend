package routes

import (
	"wow-bato-backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterBudgetItemRoutes(router *gin.RouterGroup){
	budgetItem := router.Group("/budgetItem")
	{
		budgetItem.POST("/add/:categoryID", handlers.AddNewBudgetItem)
		budgetItem.GET("/all/:categoryID", handlers.GetAllBudgetItem)
		budgetItem.GET("/:categoryID/:budgetItemID", handlers.GetSingleBudgetItem)
	}
}