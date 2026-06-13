package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"

	"github.com/aidilbaihaqi/prabodrive-be/internal/config"
	"github.com/aidilbaihaqi/prabodrive-be/internal/delivery/http/handler"
	"github.com/aidilbaihaqi/prabodrive-be/internal/delivery/http/routes"
	"github.com/aidilbaihaqi/prabodrive-be/internal/infrastructure/database"
	"github.com/aidilbaihaqi/prabodrive-be/internal/middleware"
	"github.com/aidilbaihaqi/prabodrive-be/internal/repository"
	"github.com/aidilbaihaqi/prabodrive-be/internal/services"
	"github.com/aidilbaihaqi/prabodrive-be/internal/usecase"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file, using system environment")
	}

	cfg := config.Load()
	gin.SetMode(cfg.GinMode)

	ctx := context.Background()

	db, err := database.NewPool(ctx, cfg.DB)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer db.Close()
	log.Println("database connected")

	awsCfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(cfg.AWS.Region),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AWS.AccessKeyID,
			cfg.AWS.SecretAccessKey,
			"",
		)),
	)
	if err != nil {
		log.Fatalf("aws config: %v", err)
	}

	s3Svc, err := services.NewS3Service(*cfg)
	if err != nil {
		log.Fatalf("s3 service: %v", err)
	}

	sesClient := ses.NewFromConfig(awsCfg)
	emailSvc := services.NewEmailService(sesClient, cfg.SES)

	// Repositories
	userRepo := repository.NewUserRepository(db)
	tokenRepo := repository.NewRefreshTokenRepository(db)
	docRepo := repository.NewDocumentRepository(db)
	folderRepo := repository.NewFolderRepository(db)
	shareRepo := repository.NewShareRepository(db)
	activityRepo := repository.NewActivityRepository(db)

	// Services
	quotaSvc := services.NewQuotaService(userRepo)

	// Base URL for share links
	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = "http://localhost:" + cfg.Port
	}

	// Usecases
	authUC := usecase.NewAuthUsecase(userRepo, tokenRepo, activityRepo, cfg.JWT)
	docUC := usecase.NewDocumentUsecase(docRepo, userRepo, activityRepo, s3Svc, quotaSvc)
	folderUC := usecase.NewFolderUsecase(folderRepo)
	shareUC := usecase.NewShareUsecase(shareRepo, docRepo, userRepo, activityRepo, s3Svc, emailSvc)
	adminUC := usecase.NewAdminUsecase(userRepo)

	// Handlers
	authH := handler.NewAuthHandler(authUC)
	docH := handler.NewDocumentHandler(docUC)
	folderH := handler.NewFolderHandler(folderUC)
	shareH := handler.NewShareHandler(shareUC, baseURL)
	activityH := handler.NewActivityHandler(activityRepo)
	adminH := handler.NewAdminHandler(adminUC)

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(middleware.RateLimit())
	r.Use(middleware.MaintenanceMode())
	r.Use(middleware.CORS())

	routes.Register(r, cfg.JWT.Secret, authH, docH, folderH, shareH, activityH, adminH)

	log.Printf("server starting on :%s (mode: %s)", cfg.Port, cfg.GinMode)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server: %v", err)
	}
}
