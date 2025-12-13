package service

import (
	"context"
	"fmt"
	"ingredient-recognition-backend/internal/aws"
	"ingredient-recognition-backend/internal/domain"
	"mime/multipart"
	"strings"
)

// DetectorService defines the interface for detecting ingredients from images.
type DetectorService interface {
	DetectIngredientsFromImage(ctx context.Context, file *multipart.FileHeader) ([]domain.Ingredient, error)
	DetectIngredientsFromImageWithCustomLabels(ctx context.Context, file *multipart.FileHeader) ([]domain.Ingredient, error)
}

// detectorService is a concrete implementation of the DetectorService interface.
type detectorService struct {
	awsClient *aws.AWSClient
	config    *DetectorConfig
}

// DetectorConfig holds configuration for the detector service
type DetectorConfig struct {
	ProjectARN    string
	ModelVersion  string
	MinConfidence float32
}

// NewDetectorService creates a new instance of DetectorService.
func NewDetectorService(awsClient *aws.AWSClient) DetectorService {
	return &detectorService{awsClient: awsClient}
}

// NewDetectorServiceWithCustomLabels creates a new DetectorService with custom labels configuration
func NewDetectorServiceWithCustomLabels(awsClient *aws.AWSClient, config *DetectorConfig) DetectorService {
	return &detectorService{
		awsClient: awsClient,
		config:    config,
	}
}

// DetectIngredientsFromImage reads an uploaded file and detects ingredients.
func (d *detectorService) DetectIngredientsFromImage(ctx context.Context, file *multipart.FileHeader) ([]domain.Ingredient, error) {
	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	// Read file contents
	buf := make([]byte, file.Size)
	if _, err := src.Read(buf); err != nil {
		return nil, err
	}

	// Detect ingredients from image data
	labels, err := d.awsClient.Rekognition.DetectLabels(ctx, buf)
	if err != nil {
		return nil, err
	}

	// Convert labels to ingredients
	ingredients := d.labelsToIngredients(labels)
	return ingredients, nil
}

// DetectIngredientsFromImageWithCustomLabels reads an uploaded file and detects ingredients using custom labels
func (d *detectorService) DetectIngredientsFromImageWithCustomLabels(ctx context.Context, file *multipart.FileHeader) ([]domain.Ingredient, error) {
	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	// Read file contents
	buf := make([]byte, file.Size)
	if _, err := src.Read(buf); err != nil {
		return nil, err
	}

	// Detect ingredients from image data using custom labels
	if d.config == nil {
		return nil, fmt.Errorf("custom labels configuration not set")
	}

	labels, err := d.awsClient.Rekognition.DetectCustomLabels(ctx, buf, d.config.ProjectARN, d.config.ModelVersion, d.config.MinConfidence)
	if err != nil {
		return nil, err
	}

	// Convert labels with confidence scores to ingredients
	ingredients := d.customLabelsToIngredients(labels)
	return ingredients, nil
}

// labelsToIngredients converts AWS Rekognition labels to domain Ingredient objects
func (d *detectorService) labelsToIngredients(labels []string) []domain.Ingredient {
	var ingredients []domain.Ingredient

	// Map of common food labels to ingredient names
	foodKeywords := map[string]bool{
		"apple": true, "banana": true, "orange": true, "bread": true,
		"cheese": true, "milk": true, "egg": true, "tomato": true,
		"carrot": true, "potato": true, "onion": true, "garlic": true,
		"chicken": true, "beef": true, "fish": true, "rice": true,
		"pasta": true, "vegetable": true, "fruit": true, "meat": true,
		"butter": true, "oil": true, "salt": true, "pepper": true,
	}

	for _, label := range labels {
		lowerLabel := strings.ToLower(label)

		// Check if label is a food-related keyword
		for keyword := range foodKeywords {
			if strings.Contains(lowerLabel, keyword) {
				ingredient := domain.Ingredient{
					Name:     label,
					Quantity: 1.0,
					Unit:     "unit",
				}
				ingredients = append(ingredients, ingredient)
				break
			}
		}
	}

	return ingredients
}

// customLabelsToIngredients converts custom labels with confidence scores to domain Ingredient objects
func (d *detectorService) customLabelsToIngredients(labels map[string]float32) []domain.Ingredient {
	var ingredients []domain.Ingredient

	// Filter labels by ingredient categories
	ingredientCategories := map[string]bool{
		"apple": true, "banana": true, "orange": true, "bread": true,
		"cheese": true, "milk": true, "egg": true, "tomato": true,
		"carrot": true, "potato": true, "onion": true, "garlic": true,
		"chicken": true, "beef": true, "fish": true, "rice": true,
		"pasta": true, "vegetable": true, "fruit": true, "meat": true,
		"butter": true, "oil": true, "salt": true, "pepper": true,
	}

	for label, confidence := range labels {
		lowerLabel := strings.ToLower(label)

		// Check if label is an ingredient category
		if ingredientCategories[lowerLabel] {
			ingredient := domain.Ingredient{
				Name:     label,
				Quantity: float64(confidence),
				Unit:     "confidence",
			}
			ingredients = append(ingredients, ingredient)
		}
	}

	return ingredients
}
