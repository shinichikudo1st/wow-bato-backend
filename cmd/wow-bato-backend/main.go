package main

import (
	"log"
	"os"
	database "wow-bato-backend/internal"
	"wow-bato-backend/internal/handlers"
	"wow-bato-backend/internal/routes"
	"wow-bato-backend/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

    db, err := database.ConnectDB()
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

	router := gin.Default()

	// Stores the SESSION_SECRET in a cookie
	store := cookie.NewStore([]byte(os.Getenv("SESSION_SECRET")))
	store.Options(sessions.Options{
		Path: "/",
		Domain: "localhost",
		HttpOnly: true,
		MaxAge: 0,
	})

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

	barangayService := services.NewBarangayService(db)
	barangayHandler := handlers.NewBarangayHandlers(barangayService)

	v1 := router.Group("/api/v1")
	{
		routes.RegisterUserRoute(v1)

		routes.RegisterBarangayRoute(v1, barangayHandler)

		routes.RegisterBudgetCategoryRoutes(v1)

		routes.RegisterBudgetItemRoutes(v1)

		routes.RegisterProjectRoutes(v1)

		routes.RegisterFeedbackRoutes(v1)

		routes.RegisterFeedbackReplyRoutes(v1)
	}

	router.Run(":8080")

}
