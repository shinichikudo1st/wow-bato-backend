package routes

import (
	"wow-bato-backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

func AddBarangayRoute(router *gin.Engine){
	router.POST("/addBarangay", handlers.AddBarangay)
}

func DeleteBarangayRoute(router *gin.Engine){
	router.DELETE("/deleteBarangay", handlers.DeleteBarangay)
}