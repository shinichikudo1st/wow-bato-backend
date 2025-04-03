package handlers

import (
	"net/http"
	"wow-bato-backend/internal/models"
	"wow-bato-backend/internal/services"

	"github.com/gin-gonic/gin"
)

func AddBarangay(c *gin.Context) {

	services.CheckAuthentication(c)

	var newBarangay models.AddBarangay
	services.BindJSON(c, &newBarangay)

	err := services.AddNewBarangay(newBarangay)
	services.CheckServiceError(c, err)

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Successfully Added New Barangay"})

}

func GetAllBarangay(c *gin.Context){

	services.CheckAuthentication(c)

	page := c.Query("page")
	limit := c.Query("limit")

	barangay, err := services.GetAllBarangay(limit, page)
	services.CheckServiceError(c, err)

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Successfully fetched Barangays", "data": barangay})
}

func GetSingleBarangay(c *gin.Context){

	services.CheckAuthentication(c)

	barangay_ID := c.Param("barangay_ID")

	barangay, err := services.GetSingleBarangay(barangay_ID)
	services.CheckServiceError(c, err)

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Retrieved specific barangay", "data": barangay})
}


func DeleteBarangay(c *gin.Context) {

	services.CheckAuthentication(c)

	barangay_ID := c.Param("barangay_ID")

	err := services.DeleteBarangay(barangay_ID)
	services.CheckServiceError(c, err)

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Successfully deleted the Barangay"})
}


func UpdateBarangay(c *gin.Context) {

	services.CheckAuthentication(c)

	barangay_ID := c.Param("barangay_ID")

	var barangayUpdate models.UpdateBarangay
	services.BindJSON(c, &barangayUpdate)

	err := services.UpdateBarangay(barangay_ID, barangayUpdate)
	services.CheckServiceError(c, err)

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Successfully Updated Barangay"})
}

func GetBarangayOptions(c *gin.Context){
   
    barangay, err := services.OptionBarangay()
	services.CheckServiceError(c, err)

    c.IndentedJSON(http.StatusOK, gin.H{"message": "Barangays found", "data": barangay})
}

func GetPublicBarangay(c *gin.Context){

    barangays, err := services.AllBarangaysPublic()
	services.CheckServiceError(c, err)


    c.IndentedJSON(http.StatusOK, gin.H{"message": "All barangays retrieved","data": barangays})
}
