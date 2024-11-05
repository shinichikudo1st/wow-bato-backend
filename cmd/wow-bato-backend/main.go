package main

import (
	"log"
	"os"
	"wow-bato-backend/internal/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	router := gin.Default()

	store := cookie.NewStore([]byte(os.Getenv("SESSION_SECRET")))
	router.Use(sessions.Sessions("mysession", store))

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	routes.RegisterUserRoute(router)
	routes.LoginUserRoute(router)
	routes.LogoutUserRoute(router)
	routes.CheckAuthRoute(router)

	routes.AddBarangayRoute(router)
	routes.DeleteBarangayRoute(router)
	routes.UpdateBarangayRoute(router)
	routes.GetAllBarangay(router)

	router.Run(":8080")

}
