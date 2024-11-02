package main

import (
	"wow-bato-backend/internal/routes"

	"github.com/gin-gonic/gin"
)

func main(){
	router := gin.Default()
	
	routes.RegisterUserRoute(router)
	router.Run(":8080")

}