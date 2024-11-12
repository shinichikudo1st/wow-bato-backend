package routes

import "github.com/gin-gonic/gin"

func RegisterProjectRoutes(router *gin.RouterGroup){
	project := router.Group("/project")
	{
		project.POST("/add/:categoryID")
		project.DELETE("/delete/:projectID")
	}
}