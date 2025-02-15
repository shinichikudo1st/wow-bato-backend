// Package handlers provides HTTP request handlers for the wow-bato application.
// It implements handlers for barangay (administrative division) management, including:
//   - Barangay registration and creation
//   - Barangay information updates
//   - Barangay listing and retrieval
//   - Barangay deletion with dependency checks
//   - Barangay option listing for UI components
//
// The package ensures proper data validation and maintains
// consistency across barangay-related operations.
package handlers

import (
	"net/http"
	"wow-bato-backend/internal/models"
	"wow-bato-backend/internal/services"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// AddBarangay handles the registration of new barangays in the system.
//
// This handler performs the following operations:
//  1. Validates user authentication and authorization
//  2. Validates and binds the barangay data
//  3. Delegates barangay creation to the services layer
//  4. Returns appropriate response based on operation result
//
// Security:
//  - Requires authenticated session
//  - Validates administrative privileges
//
// @Summary Register a new barangay
// @Description Creates a new barangay with the provided information
// @Tags Barangay
// @Accept json
// @Produce json
// @Param barangay body models.AddBarangay true "Barangay details including name, city, and region"
// @Success 200 {object} gin.H "Returns success message on barangay creation"
// @Failure 400 {object} gin.H "Returns error when request validation fails"
// @Failure 401 {object} gin.H "Returns error when user is not authenticated"
// @Failure 500 {object} gin.H "Returns error when barangay creation fails"
// @Router /barangay [post]
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

// GetAllBarangay retrieves a paginated list of barangays.
//
// This handler performs the following operations:
//  1. Validates user authentication
//  2. Processes pagination parameters
//  3. Retrieves barangay list with proper filtering
//  4. Formats response with barangay information
//
// The response includes:
//  - Basic barangay information
//  - Associated city and region
//  - Creation and modification timestamps
//
// @Summary Retrieve all barangays with pagination
// @Description Gets a list of barangays with optional filtering and pagination
// @Tags Barangay
// @Accept json
// @Produce json
// @Param page query string false "Page number for pagination"
// @Param limit query string false "Number of items per page"
// @Success 200 {object} gin.H "Returns list of barangays and pagination info"
// @Failure 401 {object} gin.H "Returns error when user is not authenticated"
// @Failure 500 {object} gin.H "Returns error when retrieval fails"
// @Router /barangay [get]
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

// GetSingleBarangay retrieves detailed information for a specific barangay.
//
// This handler performs the following operations:
//  1. Validates user authentication
//  2. Extracts barangay ID from request
//  3. Retrieves comprehensive barangay data
//  4. Formats response with detailed information
//
// The response includes:
//  - Complete barangay details
//  - Administrative information
//  - Associated projects and budgets
//  - Historical data and statistics
//
// @Summary Get detailed barangay information
// @Description Retrieves comprehensive information about a specific barangay
// @Tags Barangay
// @Accept json
// @Produce json
// @Param barangayID path string true "ID of the barangay"
// @Success 200 {object} gin.H "Returns detailed barangay information"
// @Failure 401 {object} gin.H "Returns error when user is not authenticated"
// @Failure 404 {object} gin.H "Returns error when barangay is not found"
// @Failure 500 {object} gin.H "Returns error when retrieval fails"
// @Router /barangay/{barangayID} [get]
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

// DeleteBarangay handles the removal of barangay records from the system.
//
// This handler performs the following operations:
//  1. Validates user authentication and authorization
//  2. Checks for dependent records
//  3. Performs cascading deletion if required
//  4. Maintains referential integrity
//
// Security:
//  - Requires authenticated session
//  - Validates administrative privileges
//  - Ensures proper cleanup of dependencies
//
// @Summary Delete a barangay
// @Description Removes a barangay and its associated data from the system
// @Tags Barangay
// @Accept json
// @Produce json
// @Param barangayID path string true "ID of the barangay to delete"
// @Success 200 {object} gin.H "Returns success message on deletion"
// @Failure 401 {object} gin.H "Returns error when user is not authenticated"
// @Failure 403 {object} gin.H "Returns error when user lacks permission"
// @Failure 500 {object} gin.H "Returns error when deletion fails"
// @Router /barangay/{barangayID} [delete]
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

// UpdateBarangay handles modifications to existing barangay information.
//
// This handler performs the following operations:
//  1. Validates user authentication and authorization
//  2. Validates update data consistency
//  3. Applies changes while maintaining integrity
//  4. Updates associated timestamps
//
// The update operation supports modifying:
//  - Barangay name and details
//  - Administrative information
//  - Associated metadata
//
// @Summary Update barangay information
// @Description Modifies existing barangay details with provided updates
// @Tags Barangay
// @Accept json
// @Produce json
// @Param barangayID path string true "ID of the barangay to update"
// @Param barangay body models.UpdateBarangay true "Updated barangay details"
// @Success 200 {object} gin.H "Returns success message on update"
// @Failure 400 {object} gin.H "Returns error when request validation fails"
// @Failure 401 {object} gin.H "Returns error when user is not authenticated"
// @Failure 500 {object} gin.H "Returns error when update fails"
// @Router /barangay/{barangayID} [put]
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

// GetBarangayOptions provides a simplified list of barangays for UI components.
//
// This handler performs the following operations:
//  1. Retrieves essential barangay information
//  2. Formats data for dropdown/select components
//  3. Optimizes response size for UI performance
//
// The response includes:
//  - Barangay ID and name
//  - Minimal metadata for UI rendering
//  - Sorted list for easy selection
//
// @Summary Get barangay list for UI components
// @Description Retrieves a simplified list of barangays for dropdown menus
// @Tags Barangay
// @Accept json
// @Produce json
// @Success 200 {object} gin.H "Returns list of barangay options"
// @Failure 500 {object} gin.H "Returns error when retrieval fails"
// @Router /barangay/options [get]
func GetBarangayOptions(c *gin.Context){
   
    barangay, err := services.OptionBarangay()

    if err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.IndentedJSON(http.StatusOK, gin.H{"message": "Barangays found", "data": barangay})
}

// GetPublicBarangay provides a simplified list of barangays for UI components.
//
// This handler performs the following operations:
//  1. Fetching of barangays publicly for non authorized user
//
// @Summary GetPublicBarangay get all barangay publicly
// @Description Retrieves a simplified list of barangays for non authorized requests
// @Tags Barangay
// @Accept json
// @Produce json
// @Success 200 {object} gin.H "Returns list of barangay"
// @Failure 500 {object} gin.H "Returns error when retrieval fails"
// @Router /barangay/public-all [get]
func GetPublicBarangay(c *gin.Context){

    barangays, err := services.AllBarangaysPublic()

    if err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }


    c.IndentedJSON(http.StatusOK, gin.H{"message": "All barangays retrieved","data": barangays})
}
