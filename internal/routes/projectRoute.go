package routes

import (
	"wow-bato-backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterProjectRoutes(router *gin.RouterGroup){
	project := router.Group("/project")
	{
		project.POST("/add/:categoryID", handlers.AddNewProject)
		project.DELETE("/delete/:projectID", handlers.DeleteProject)
		project.GET("/all", handlers.GetAllProjects)
		project.PATCH("/update-status/:projectID", handlers.UpdateProjectStatus)
	}
}