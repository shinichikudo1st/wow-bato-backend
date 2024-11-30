package routes

import (
	"wow-bato-backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterBarangayRoute(router *gin.RouterGroup){
	barangay := router.Group("/barangay")
	{
		barangay.POST("/add", handlers.AddBarangay)
		barangay.DELETE("/delete/:barangay_ID", handlers.DeleteBarangay)
		barangay.PUT("/update/:barangay_ID", handlers.UpdateBarangay)
		barangay.GET("/all", handlers.GetAllBarangay)
		barangay.GET("/single/:barangay_ID", handlers.GetSingleBarangay)
        barangay.GET("/options", handlers.GetBarangayOptions)
	}
}
