package routes

import (
	"wow-bato-backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterBarangayRoute registers barangay-related routes on the provided router group.
//
// This route contains:
// - POST /barangay/add: Handles adding a new barangay.
// - DELETE /barangay/delete/:barangay_ID: Handles deleting a barangay.
// - PUT /barangay/update/:barangay_ID: Handles updating a barangay.
// - GET /barangay/all: Handles retrieving a list of barangays.
// - GET /barangay/single/:barangay_ID: Handles retrieving a single barangay.
// - GET /barangay/options: Handles retrieving barangay used for dropdown selection
func RegisterBarangayRoute(router *gin.RouterGroup) {
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
