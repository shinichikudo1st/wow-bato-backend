package handlers

import (
	"net/http"
	"wow-bato-backend/internal/models"
	"wow-bato-backend/internal/services"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

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

func GetAllBarangay(c *gin.Context){
	session := sessions.Default(c)

	if session.Get("authenticated") != true {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Access Denied"})
		return
	}

	page := c.Query("page")
	limit := c.Query("limit")

	barangay, err := services.GetAllBarangay(page, limit)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Successfully fetched Barangays", "data": barangay})
}

func DeleteBarangay(c *gin.Context) {
	session := sessions.Default(c)

	if session.Get("authenticated") != true {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Access Denied"})
		return
	}

	var barangayToDelete models.DeleteBarangay

	if err := c.ShouldBindJSON(&barangayToDelete); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := services.DeleteBarangay(barangayToDelete)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Successfully deleted the Barangay"})
}

func UpdateBarangay(c *gin.Context) {
	session := sessions.Default(c)

	if session.Get("authenticated") != true {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized: Access Denied"})
		return
	}

	var updateBarangay models.UpdateBarangay

	if err := c.ShouldBindJSON(&updateBarangay); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.UpdateBarangay(updateBarangay); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Successfully Updated Barangay"})
}
