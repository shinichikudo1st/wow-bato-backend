package main

import (
	"fmt"
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
	"gorm.io/gorm"
)

type App struct {
	DB                     *gorm.DB
    BarangayHandlers       *handlers.BarangayHandlers
    UserHandlers           *handlers.UserHandlers
    BudgetItemHandlers     *handlers.BudgetItemHandlers
    FeedbackHandlers       *handlers.FeedbackHandlers
    FeedbackReplyHandlers  *handlers.FeedbackReplyHandlers
    BudgetCategoryHandlers *handlers.BudgetCategoryHandlers
    ProjectHandlers        *handlers.ProjectHandlers
}

func NewApp() (*App, error) {
    db, err := database.ConnectDB()
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }

    barangayService := services.NewBarangayService(db)
    userService := services.NewUserService(db)
    budgetItemService := services.NewBudgetItemService(db)
    feedbackService := services.NewFeedbackService(db)
    feedbackReplyService := services.NewFeedbackReplyService(db)
    budgetCategoryService := services.NewBudgetCategoryService(db)
    projectService := services.NewProjectService(db)

    return &App{
        DB:                     db,
        BarangayHandlers:       handlers.NewBarangayHandlers(barangayService),
        UserHandlers:           handlers.NewUserHandlers(userService),
        BudgetItemHandlers:     handlers.NewBudgetItemHandlers(budgetItemService),
        FeedbackHandlers:       handlers.NewFeedbackHandlers(feedbackService),
        FeedbackReplyHandlers:  handlers.NewFeedbackReplyHandlers(feedbackReplyService),
        BudgetCategoryHandlers: handlers.NewBudgetCategoryHandlers(budgetCategoryService),
        ProjectHandlers:        handlers.NewProjectHandlers(projectService, budgetCategoryService),
    }, nil
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

    app, err := NewApp()
    if err != nil {
        log.Fatalf("Failed to initialize app: %v", err)
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


	v1 := router.Group("/api/v1")
	{
		routes.RegisterUserRoute(v1, app.UserHandlers)
		routes.RegisterBarangayRoute(v1, app.BarangayHandlers)
		routes.RegisterBudgetCategoryRoutes(v1, app.BudgetCategoryHandlers)
		routes.RegisterBudgetItemRoutes(v1, app.BudgetItemHandlers)
		routes.RegisterProjectRoutes(v1, app.ProjectHandlers)
		routes.RegisterFeedbackRoutes(v1, app.FeedbackHandlers)
		routes.RegisterFeedbackReplyRoutes(v1, app.FeedbackReplyHandlers)
	}

	router.Run(":8080")

}
