package main

import (
	"context"
	"ingredient-recognition-backend/internal/aws"
	"ingredient-recognition-backend/internal/config"
	"ingredient-recognition-backend/internal/handler"
	"ingredient-recognition-backend/internal/middleware"
	"ingredient-recognition-backend/internal/repository"
	"ingredient-recognition-backend/internal/service"
	"ingredient-recognition-backend/pkg/logger"
	"log"
	"time"

	// "time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// Initialize structured logger with zap
	if err := logger.InitializeGlobalLogger("logs/app.log", true); err != nil {
		log.Fatalf("could not initialize logger: %v", err)
	}
	defer func() {
		if l := logger.GetLogger(); l != nil {
			l.Sync()
		}
	}()

	ctx := context.Background()
	logger.Info(ctx, "Application starting")

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal(ctx, "Failed to load configuration", err)
	}

	// Initialize AWS client
	awsClient, err := aws.NewAWSClient(context.TODO(), cfg.AWSRegion, cfg.AWSBucket)
	if err != nil {
		logger.Fatal(ctx, "Failed to initialize AWS client", err, zap.String("region", cfg.AWSRegion))
	}
	logger.Info(ctx, "AWS client initialized", zap.String("region", cfg.AWSRegion))

	// Initialize user repository
	userRepo := repository.NewUserRepository(awsClient.DynamoDB)

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
			ModelArn:      cfg.RekognitionModelARN,
			ProjectARN:    cfg.RekognitionProjectARN,
			ModelVersion:  cfg.RekognitionModelVersion,
			MinConfidence: cfg.RekognitionMinConfidence,
		}
		customDetectorService = service.NewDetectorServiceWithCustomLabels(awsClient, customConfig)
		ingredientHandler = handler.NewIngredientHandler(customDetectorService)
	}

	// Initialize recipe service with Bedrock
	recipeService := service.NewRecipeService(awsClient.BedrockRuntime, cfg.BedrockModelID)
	recipeHandler := handler.NewRecipeHandler(recipeService)

	// Create Gin router
	router := gin.Default()

	// Add logging and error handling middleware
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.ErrorHandlingMiddleware())

	// Public routes (no auth required)
	router.POST("/auth/register", authHandler.Register)
	router.POST("/auth/login", authHandler.Login)

	router.GET("/health", func(c *gin.Context) {
		logger.Debug(c.Request.Context(), "Health check requested")
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Protected routes (auth required)
	protected := router.Group("/api")

	routeVersion := protected.Group("/v1")
	routeVersion.Use(middleware.AuthMiddleware(authService))
	routeVersion.POST("/detect", ingredientHandler.DetectIngredientsWithCustomLabels)
	routeVersion.POST("/recipes/recommend", recipeHandler.RecommendRecipes)

	// Start the server
	logger.Info(ctx, "Starting server", zap.String("address", cfg.ServerAddress))
	if err := router.Run(cfg.ServerAddress); err != nil {
		logger.Fatal(ctx, "Failed to start server", err, zap.String("address", cfg.ServerAddress))
	}
}
