package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/yourname/yourapp/internal/config"
	"github.com/yourname/yourapp/internal/delivery/http/handler"
	"github.com/yourname/yourapp/internal/delivery/http/routes"
	"github.com/yourname/yourapp/internal/infrastructure/database"
	"github.com/yourname/yourapp/internal/middleware"
	"github.com/yourname/yourapp/internal/repository"
	"github.com/yourname/yourapp/internal/usecase/user"
)

func main() {
	// 1. Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment")
	}

	// 2. Load configuration
	cfg := config.Load()

	// 3. Set Gin mode
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 4. Initialize infrastructure
	db, err := database.NewPostgres(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("✓ Database connected")

	// 5. Initialize repositories
	userRepo := repository.NewUserRepository(db)

	// 6. Initialize usecases
	createUserUC := user.NewCreateUserUsecase(userRepo)
	getUserUC := user.NewGetUserUsecase(userRepo)
	listUsersUC := user.NewListUsersUsecase(userRepo)
	updateUserUC := user.NewUpdateUserUsecase(userRepo)
	deleteUserUC := user.NewDeleteUserUsecase(userRepo)

	// 7. Initialize handlers
	userHandler := handler.NewUserHandler(
		createUserUC,
		getUserUC,
		listUsersUC,
		updateUserUC,
		deleteUserUC,
	)

	// 8. Setup router
	router := gin.Default()

	// 9. Register middleware
	router.Use(middleware.CORS())
	router.Use(middleware.RequestID())
	router.Use(middleware.Logger())

	// 10. Register routes
	routes.RegisterHealthRoutes(router)
	routes.RegisterUserRoutes(router, userHandler, cfg.JWT.AccessSecret)

	// 11. Start server
	port := cfg.AppPort
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("📝 Environment: %s", cfg.AppEnv)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
		os.Exit(1)
	}
}
