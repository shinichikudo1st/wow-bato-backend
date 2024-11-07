package routes

import (
	"wow-bato-backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

func AddBudgetCategoryRoute(router *gin.Engine) {
	router.POST("/addBudgetCategory", handlers.AddBudgetCategory)
}

func DeleteBudgetCategoryRoute(router *gin.Engine) {
	router.DELETE("/deleteBudgetCategory/:budgetID", handlers.DeleteBudgetCategory)
}