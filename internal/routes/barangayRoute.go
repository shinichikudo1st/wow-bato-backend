package routes

import "github.com/gin-gonic/gin"

func AddBarangayRoute(router *gin.Engine){
	router.POST("/addBarangay")
}