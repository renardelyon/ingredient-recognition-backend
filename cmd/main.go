package main

import (
	"context"
	"ingredient-recognition-backend/internal/aws"
	"ingredient-recognition-backend/internal/config"
	"ingredient-recognition-backend/internal/handler"
	"ingredient-recognition-backend/internal/service"
	"log"

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

	// Initialize services and handlers
	detectorService := service.NewDetectorService(awsClient)
	ingredientHandler := handler.NewIngredientHandler(detectorService)

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

	// Set up routes
	router.POST("/detect", ingredientHandler.DetectIngredients)
	router.POST("/detect-custom", ingredientHandler.DetectIngredientsWithCustomLabels)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Start the server
	log.Printf("Starting server on %s", cfg.ServerAddress)
	if err := router.Run(cfg.ServerAddress); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
