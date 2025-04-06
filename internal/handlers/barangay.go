package handlers

import (
	"net/http"
	"wow-bato-backend/internal/models"
	"wow-bato-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type BarangayHandlers struct {
	svc *services.BarangayService
}

func NewBarangayHandlers(svc *services.BarangayService) *BarangayHandlers {
	return &BarangayHandlers{svc: svc}
}

func (h *BarangayHandlers) AddBarangay(c *gin.Context) {

	services.CheckAuthentication(c)

	var newBarangay models.AddBarangay
	services.BindJSON(c, &newBarangay)

	err := h.svc.AddNewBarangay(newBarangay)
	services.CheckServiceError(c, err)

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Successfully Added New Barangay"})

}

func (h *BarangayHandlers) GetAllBarangay(c *gin.Context){

	services.CheckAuthentication(c)

	page := c.Query("page")
	limit := c.Query("limit")

	barangay, err := h.svc.GetAllBarangay(limit, page)
	services.CheckServiceError(c, err)

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Successfully fetched Barangays", "data": barangay})
}

func (h *BarangayHandlers) GetSingleBarangay(c *gin.Context){

	services.CheckAuthentication(c)

	barangay_ID := c.Param("barangay_ID")

	barangay, err := h.svc.GetSingleBarangay(barangay_ID)
	services.CheckServiceError(c, err)

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Retrieved specific barangay", "data": barangay})
}


func (h *BarangayHandlers) DeleteBarangay(c *gin.Context) {

	services.CheckAuthentication(c)

	barangay_ID := c.Param("barangay_ID")

	err := h.svc.DeleteBarangay(barangay_ID)
	services.CheckServiceError(c, err)

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Successfully deleted the Barangay"})
}


func (h *BarangayHandlers) UpdateBarangay(c *gin.Context) {

	services.CheckAuthentication(c)

	barangay_ID := c.Param("barangay_ID")

	var barangayUpdate models.UpdateBarangay
	services.BindJSON(c, &barangayUpdate)

	err := h.svc.UpdateBarangay(barangay_ID, barangayUpdate)
	services.CheckServiceError(c, err)

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Successfully Updated Barangay"})
}

func (h *BarangayHandlers) GetBarangayOptions(c *gin.Context){
   
    barangay, err := h.svc.OptionBarangay()
	services.CheckServiceError(c, err)

    c.IndentedJSON(http.StatusOK, gin.H{"message": "Barangays found", "data": barangay})
}

func (h *BarangayHandlers) GetPublicBarangay(c *gin.Context){

    barangays, err := h.svc.AllBarangaysPublic()
	services.CheckServiceError(c, err)


    c.IndentedJSON(http.StatusOK, gin.H{"message": "All barangays retrieved","data": barangays})
}
