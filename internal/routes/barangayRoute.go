package routes

import (
	"wow-bato-backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

func AddBarangayRoute(router *gin.Engine) {
	router.POST("/addBarangay", handlers.AddBarangay)
}

func DeleteBarangayRoute(router *gin.Engine) {
	router.DELETE("/deleteBarangay", handlers.DeleteBarangay)
}

func UpdateBarangayRoute(router *gin.Engine) {
	router.PUT("/updateBarangay", handlers.UpdateBarangay)
}

func GetAllBarangay(router *gin.Engine){
	router.GET("/getAllBarangay", handlers.GetAllBarangay)
}

func GetSingleBarangay(router *gin.Engine){
	router.GET("/getSingleBarangay", handlers.GetSingleBarangay)
}
