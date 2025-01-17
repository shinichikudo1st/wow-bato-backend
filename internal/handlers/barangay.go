package handlers

import (
	"net/http"
	"wow-bato-backend/internal/models"
	"wow-bato-backend/internal/services"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Handler for adding barangay
// @Summary Add barangay
// @Tags Barangay
// @Accept json
// @Produce json
// @Param barangay body models.AddBarangay true "Barangay details"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
func AddBarangay(c *gin.Context) {
	session := sessions.Default(c)

	if session.Get("authenticated") != true {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Access Denied"})
		return
	}

	var newBarangay models.AddBarangay

	if err := c.ShouldBindJSON(&newBarangay); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := services.AddNewBarangay(newBarangay)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Successfully Added New Barangay"})

}

// Handler for getting barangay
// @Summary Get barangay
// @Tags Barangay
// @Accept json no body
// @Produce json
// @Success 200 {object} gin.H
// @Failure 500 {object} gin.H
func GetAllBarangay(c *gin.Context){
	session := sessions.Default(c)

	if session.Get("authenticated") != true {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Access Denied"})
		return
	}

	page := c.Query("page")
	limit := c.Query("limit")

	barangay, err := services.GetAllBarangay(limit, page)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Successfully fetched Barangays", "data": barangay})
}

// Handler for getting barangay
// @Summary Get barangay
// @Tags Barangay
// @Accept json no body
// @Produce json
// @Success 200 {object} gin.H
// @Failure 500 {object} gin.H
func GetSingleBarangay(c *gin.Context){
	session := sessions.Default(c)

	if session.Get("authenticated") != true {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Access Denied"})
		return
	}

	barangay_ID := c.Param("barangay_ID")

	barangay, err := services.GetSingleBarangay(barangay_ID)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Retrieved specific barangay", "data": barangay})
}

// Handler for deleting barangay
// @Summary Delete barangay
// @Tags Barangay
// @Accept json no body
// @Produce json
// @Success 200 {object} gin.H
// @Failure 500 {object} gin.H
func DeleteBarangay(c *gin.Context) {
	session := sessions.Default(c)

	if session.Get("authenticated") != true {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Access Denied"})
		return
	}

	barangay_ID := c.Param("barangay_ID")

	err := services.DeleteBarangay(barangay_ID)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Successfully deleted the Barangay"})
}

// Handler for updating barangay
// @Summary Update barangay
// @Tags Barangay
// @Accept json
// @Produce json
// @Param barangay body models.UpdateBarangay true "Barangay details"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
func UpdateBarangay(c *gin.Context) {
	session := sessions.Default(c)

	if session.Get("authenticated") != true {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized: Access Denied"})
		return
	}

	barangay_ID := c.Param("barangay_ID")

	var barangayUpdate models.UpdateBarangay

	if err := c.ShouldBindJSON(&barangayUpdate); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := services.UpdateBarangay(barangay_ID, barangayUpdate)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Successfully Updated Barangay"})
}

// Handler for getting barangay options
// @Summary Get barangay for dropdown selection on the client side
// @Tags Barangay
// @Accept json no body
// @Produce json
// @Success 200 {object} gin.H
// @Failure 500 {object} gin.H
func GetBarangayOptions(c *gin.Context){
   
    barangay, err := services.OptionBarangay()

    if err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.IndentedJSON(http.StatusOK, gin.H{"message": "Barangays found", "data": barangay})
}
