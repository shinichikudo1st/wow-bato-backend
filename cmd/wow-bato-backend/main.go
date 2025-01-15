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
	// Check if .env file exists
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	router := gin.Default()

	// Stores the SESSION_SECRET in a cookie
	store := cookie.NewStore([]byte(os.Getenv("SESSION_SECRET")))

	// Creates a session using the stored SESSION_SECRET key
	router.Use(sessions.Sessions("mysession", store))

	// Lists the allowed methods, headers, origin, and allow credentials
	router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// API versioning for flexibility and scalability
	v1 := router.Group("/api/v1")
	{
		// User Routes API Version 1
		routes.RegisterUserRoute(v1)

		// Barangay Routes API Version 1
		routes.RegisterBarangayRoute(v1)

		// Budget Category Routes API Version 1
		routes.RegisterBudgetCategoryRoutes(v1)

		// Budget Item Routes API Version 1
		routes.RegisterBudgetItemRoutes(v1)

		// Project Routes API Version 1
		routes.RegisterProjectRoutes(v1)

		// Feedback Routes API Version 1
		routes.RegisterFeedbackRoutes(v1)

		// Feedback Replies Routes API Version 1
		routes.RegisterFeedbackReplyRoutes(v1)
	}


	// Server will run on port 8080
	router.Run(":8080")

}
