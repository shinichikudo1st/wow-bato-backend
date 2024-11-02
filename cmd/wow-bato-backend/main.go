package main

import (
	"wow-bato-backend/internal/routes"

	"github.com/gin-gonic/gin"
)

func main(){
	router := gin.Default()
	
	routes.RegisterUserRoute(router)
	routes.LoginUserRoute(router)

	router.Run(":8080")

}