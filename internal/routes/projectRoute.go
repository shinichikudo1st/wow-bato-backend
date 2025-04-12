package routes

import (
	"wow-bato-backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterProjectRoutes(router *gin.RouterGroup, handlers *handlers.ProjectHandlers) {
	project := router.Group("/project")
	{
		project.POST("/add/:categoryID", handlers.AddNewProject)
		project.DELETE("/delete/:projectID", handlers.DeleteProject)
		project.GET("/all/:categoryID", handlers.GetAllProjects)
		project.PATCH("/update-status/:projectID", handlers.UpdateProjectStatus)
		project.GET("/specific-project/:projectID", handlers.GetSingleProject)
	}
}
