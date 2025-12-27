package service

import (
	"context"
	"encoding/json"
	"fmt"
	"ingredient-recognition-backend/internal/model"
	"ingredient-recognition-backend/internal/request"
	"ingredient-recognition-backend/pkg/logger"
	"ingredient-recognition-backend/pkg/utils"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"go.uber.org/zap"
)

// RecipeService defines methods for recipe recommendations
type RecipeService interface {
	RecommendRecipes(ctx context.Context, ingredients []string) (*model.RecipeRecommendation, error)
}

// recipeService is a concrete implementation of RecipeService
type recipeService struct {
	bedrockClient *bedrockruntime.Client
	modelID       string
}

// NewRecipeService creates a new recipe service
func NewRecipeService(bedrockClient *bedrockruntime.Client, modelID string) RecipeService {
	return &recipeService{
		bedrockClient: bedrockClient,
		modelID:       modelID,
	}
}

// RecommendRecipes generates recipe recommendations based on ingredients
func (r *recipeService) RecommendRecipes(ctx context.Context, ingredients []string) (*model.RecipeRecommendation, error) {
	logger.Info(ctx, "Starting recipe recommendation", zap.Int("ingredient_count", len(ingredients)), zap.String("ingredients", strings.Join(ingredients, ", ")))

	if len(ingredients) == 0 {
		logger.Warn(ctx, "Recipe recommendation requested with no ingredients")
		return nil, fmt.Errorf("at least one ingredient is required")
	}

	// Build the prompt for Claude
	prompt := buildRecipePrompt(ingredients)
	logger.Debug(ctx, "Generated prompt for Bedrock", zap.Int("prompt_length", len(prompt)))

	// Call Bedrock with Claude
	response, err := r.callBedrock(ctx, prompt)
	if err != nil {
		logger.Error(ctx, "Failed to call Bedrock API", err, zap.String("model_id", r.modelID))
		return nil, fmt.Errorf("failed to call Bedrock: %w", err)
	}

	// Parse the response
	recommendation, err := parseRecipeResponse(response, ingredients)
	if err != nil {
		logger.Error(ctx, "Failed to parse recipe response", err)
		return nil, fmt.Errorf("failed to parse recipe response: %w", err)
	}

	return recommendation, nil
}

// callBedrock invokes the Bedrock API with the given prompt
func (r *recipeService) callBedrock(ctx context.Context, prompt string) (string, error) {
	logger.Debug(ctx, "Calling Bedrock API", zap.String("model_id", r.modelID))

	// Prepare the request payload for Claude model
	payload := request.BedrockModelConfig{
		AnthropicVersion: "bedrock-2023-06-01",
		MaxTokens:        2048,
		Messages:         []request.Message{{Role: "user", Content: prompt}},
	}

	reqBody, err := json.Marshal(payload)
	if err != nil {
		logger.Error(ctx, "Failed to marshal Bedrock request payload", err)
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Invoke the model
	logger.Debug(ctx, "Invoking Bedrock model", zap.Int("payload_size", len(reqBody)))
	output, err := r.bedrockClient.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(r.modelID),
		ContentType: aws.String("application/json"),
		Accept:      aws.String("application/json"),
		Body:        reqBody,
	})
	if err != nil {
		logger.Error(ctx, "Bedrock model invocation failed", err, zap.String("model_id", r.modelID))
		return "", fmt.Errorf("failed to invoke model: %w", err)
	}

	// Parse the response
	var result map[string]any
	if err := json.Unmarshal(output.Body, &result); err != nil {
		logger.Error(ctx, "Failed to parse Bedrock model response", err)
		return "", fmt.Errorf("failed to parse model response: %w", err)
	}

	// Extract text from the response
	if content, ok := result["content"].([]any); ok && len(content) > 0 {
		if textBlock, ok := content[0].(map[string]any); ok {
			if text, ok := textBlock["text"].(string); ok {
				logger.Debug(ctx, "Successfully extracted text from Bedrock response", zap.Int("text_length", len(text)))
				return text, nil
			}
		}
	}

	logger.Error(ctx, "Unexpected response format from Bedrock", nil)
	return "", fmt.Errorf("unexpected response format from Bedrock")
}

// buildRecipePrompt creates a prompt for recipe generation
func buildRecipePrompt(ingredients []string) string {
	var ingredientList strings.Builder
	for i, ing := range ingredients {
		if i > 0 {
			ingredientList.WriteString(", ")
		}
		ingredientList.WriteString(ing)
	}

	return fmt.Sprintf(`Based on the following ingredients: %s

Please recommend 3-5 recipes that can be made with these ingredients. For each recipe, provide:
1. Recipe name
2. Cuisine type
3. Cooking time (in minutes)
4. Difficulty level (Easy, Medium, Hard)
5. List of ingredients needed
6. Step-by-step cooking instructions
7. Nutritional information (brief)
8. Cooking tips

Format your response as a JSON object with the following structure:
{
  "recipes": [
    {
      "name": "Recipe Name",
      "cuisine": "Cuisine Type",
      "cooking_time": "30 minutes",
      "difficulty": "Easy",
      "ingredients": ["ingredient 1", "ingredient 2"],
      "instructions": ["step 1", "step 2"],
      "nutrition": "brief nutrition info",
      "tips": "cooking tips"
    }
  ]
}

Make sure the JSON is valid and properly formatted.
Do not include any markdown formatting, explanation, or text outside the JSON object.`, ingredientList.String())
}

// parseRecipeResponse parses the Bedrock response into RecipeRecommendation
func parseRecipeResponse(responseText string, ingredients []string) (*model.RecipeRecommendation, error) {
	jsonStr, err := utils.ExtractJSONFromString(responseText)
	if err != nil {
		return nil, err
	}

	// Parse the JSON
	var recipesData model.RecipeRecommendation
	if err := json.Unmarshal([]byte(jsonStr), &recipesData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal recipes JSON: %w", err)
	}

	recipesData.IngredientCount = len(ingredients)
	recipesData.UsedIngredients = ingredients
	recipesData.GeneratedAt = fmt.Sprintf("%v", time.Now().Format(time.RFC3339))

	recipesData.TotalRecipes = len(recipesData.Recipes)
	return &recipesData, nil
}
