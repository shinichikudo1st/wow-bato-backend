package routes

import (
	"wow-bato-backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterProjectRoutes registers project-related routes on the provided router group.
//
// This route contains:
// - POST /project/add/:categoryID: Handles adding a new project.
// - DELETE /project/delete/:projectID: Handles deleting a project.
// - GET /project/all/:categoryID: Handles retrieving a list of projects.
// - PATCH /project/update-status/:projectID: Handles updating the status of a project.
func RegisterProjectRoutes(router *gin.RouterGroup) {
	project := router.Group("/project")
	{
		project.POST("/add/:categoryID", handlers.AddNewProject)
		project.DELETE("/delete/:projectID", handlers.DeleteProject)
		project.GET("/all/:categoryID", handlers.GetAllProjects)
		project.PATCH("/update-status/:projectID", handlers.UpdateProjectStatus)
	}
}
