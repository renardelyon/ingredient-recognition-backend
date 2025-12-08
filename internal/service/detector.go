package service

import (
	"context"
	"ingredient-recognition-backend/internal/aws"
	"ingredient-recognition-backend/internal/domain"
	"mime/multipart"
	"strings"
)

// DetectorService defines the interface for detecting ingredients from images.
type DetectorService interface {
	DetectIngredients(ctx context.Context, imageData []byte) ([]domain.Ingredient, error)
	DetectIngredientsFromImage(ctx context.Context, file *multipart.FileHeader) ([]domain.Ingredient, error)
}

// detectorService is a concrete implementation of the DetectorService interface.
type detectorService struct {
	awsClient *aws.AWSClient
}

// NewDetectorService creates a new instance of DetectorService.
func NewDetectorService(awsClient *aws.AWSClient) DetectorService {
	return &detectorService{awsClient: awsClient}
}

// DetectIngredients processes the given image data using AWS Rekognition and returns detected ingredients.
func (d *detectorService) DetectIngredients(ctx context.Context, imageData []byte) ([]domain.Ingredient, error) {
	// Use AWS Rekognition to detect labels
	labels, err := d.awsClient.Rekognition.DetectLabels(ctx, imageData)
	if err != nil {
		return nil, err
	}

	// Convert labels to ingredients
	ingredients := d.labelsToIngredients(labels)
	return ingredients, nil
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
	return d.DetectIngredients(ctx, buf)
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
