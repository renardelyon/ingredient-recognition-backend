package main

import (
	"context"
	"ingredient-recognition-backend/internal/aws"
	"ingredient-recognition-backend/internal/config"
	"ingredient-recognition-backend/internal/handler"
	"ingredient-recognition-backend/internal/middleware"
	"ingredient-recognition-backend/internal/repository"
	"ingredient-recognition-backend/internal/service"
	"log"
	"time"

	// "time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	// Initialize AWS client
	awsClient, err := aws.NewAWSClient(context.TODO(), cfg.AWSRegion, cfg.AWSBucket)
	if err != nil {
		log.Fatalf("could not initialize AWS client: %v", err)
	}

	// Initialize user repository
	userRepo := repository.NewUserRepository(awsClient.DynamoDB, cfg.DynamoDBTable)

	// Initialize auth service
	authService := service.NewAuthService(userRepo, cfg.JWTSecret, time.Duration(cfg.JWTExpiry)*time.Hour)

	// Initialize services and handlers
	detectorService := service.NewDetectorService(awsClient)
	ingredientHandler := handler.NewIngredientHandler(detectorService)
	authHandler := handler.NewAuthHandler(authService)

	// Initialize custom labels service if configuration is available
	var customDetectorService service.DetectorService
	if cfg.RekognitionProjectARN != "" && cfg.RekognitionModelVersion != "" {
		customConfig := &service.DetectorConfig{
			ProjectARN:    cfg.RekognitionProjectARN,
			ModelVersion:  cfg.RekognitionModelVersion,
			MinConfidence: cfg.RekognitionMinConfidence,
		}
		customDetectorService = service.NewDetectorServiceWithCustomLabels(awsClient, customConfig)
		ingredientHandler = handler.NewIngredientHandler(customDetectorService)
	}

	// Create Gin router
	router := gin.Default()

	// Public routes (no auth required)
	router.POST("/auth/register", authHandler.Register)
	router.POST("/auth/login", authHandler.Login)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Protected routes (auth required)
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware(authService))
	protected.POST("/detect", ingredientHandler.DetectIngredients)
	protected.POST("/detect-custom", ingredientHandler.DetectIngredientsWithCustomLabels)

	// Start the server
	log.Printf("Starting server on %s", cfg.ServerAddress)
	if err := router.Run(cfg.ServerAddress); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
