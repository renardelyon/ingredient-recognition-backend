package repository

import (
	"context"
	"fmt"
	"ingredient-recognition-backend/internal/domain"
	"ingredient-recognition-backend/pkg/logger"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"go.uber.org/zap"
)

// RecipeRepository is a DynamoDB implementation for recipe storage
type RecipeRepository struct {
	client    *dynamodb.Client
	tableName string
}

// NewRecipeRepository creates a new DynamoDB recipe repository
func NewRecipeRepository(client *dynamodb.Client) *RecipeRepository {
	return &RecipeRepository{
		client:    client,
		tableName: "Recipes",
	}
}

// Save creates a new saved recipe in DynamoDB
func (r *RecipeRepository) Save(ctx context.Context, recipe *domain.SavedRecipe) error {
	logger.Debug(ctx, "Saving recipe to DynamoDB", zap.String("recipe_id", recipe.ID), zap.String("user_id", recipe.UserID))

	// Marshal recipe to DynamoDB item
	item, err := attributevalue.MarshalMap(recipe)
	if err != nil {
		logger.Error(ctx, "Failed to marshal recipe", err, zap.String("recipe_id", recipe.ID))
		return fmt.Errorf("failed to marshal recipe: %w", err)
	}

	// Put item in DynamoDB
	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})

	if err != nil {
		logger.Error(ctx, "Failed to save recipe to DynamoDB", err, zap.String("recipe_id", recipe.ID))
		return fmt.Errorf("failed to save recipe: %w", err)
	}

	logger.Info(ctx, "Recipe saved successfully", zap.String("recipe_id", recipe.ID), zap.String("user_id", recipe.UserID))
	return nil
}

// GetByID retrieves a saved recipe by ID from DynamoDB
func (r *RecipeRepository) GetByID(ctx context.Context, id string) (*domain.SavedRecipe, error) {
	logger.Debug(ctx, "Getting recipe by ID", zap.String("recipe_id", id))

	result, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("id = :id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":id": &types.AttributeValueMemberS{Value: id},
		},
		Limit: aws.Int32(1),
	})

	if err != nil {
		logger.Error(ctx, "DynamoDB Query failed", err, zap.String("recipe_id", id))
		return nil, fmt.Errorf("failed to get recipe: %w", err)
	}

	if result.Count == 0 {
		return nil, domain.ErrRecipeNotFound
	}

	var recipe domain.SavedRecipe
	err = attributevalue.UnmarshalMap(result.Items[0], &recipe)
	if err != nil {
		logger.Error(ctx, "Failed to unmarshal recipe", err, zap.String("recipe_id", id))
		return nil, fmt.Errorf("failed to unmarshal recipe: %w", err)
	}

	return &recipe, nil
}

// GetByUserID retrieves all saved recipes for a user from DynamoDB
func (r *RecipeRepository) GetByUserID(ctx context.Context, userID string) ([]*domain.SavedRecipe, error) {
	logger.Debug(ctx, "Getting recipes by user ID", zap.String("user_id", userID))

	result, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String("UserIdIndex"),
		KeyConditionExpression: aws.String("user_id = :user_id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":user_id": &types.AttributeValueMemberS{Value: userID},
		},
	})

	if err != nil {
		logger.Error(ctx, "DynamoDB Query failed", err, zap.String("user_id", userID))
		return nil, fmt.Errorf("failed to get recipes: %w", err)
	}

	recipes := make([]*domain.SavedRecipe, 0, result.Count)
	for _, item := range result.Items {
		var recipe domain.SavedRecipe
		err = attributevalue.UnmarshalMap(item, &recipe)
		if err != nil {
			logger.Error(ctx, "Failed to unmarshal recipe", err)
			continue
		}
		recipes = append(recipes, &recipe)
	}

	logger.Info(ctx, "Retrieved recipes for user", zap.String("user_id", userID), zap.Int("count", len(recipes)))
	return recipes, nil
}

// Delete deletes a saved recipe from DynamoDB
func (r *RecipeRepository) Delete(ctx context.Context, id string, userID string) error {
	logger.Debug(ctx, "Deleting recipe", zap.String("recipe_id", id), zap.String("user_id", userID))

	// First verify the recipe belongs to the user
	recipe, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if recipe.UserID != userID {
		logger.Warn(ctx, "User attempted to delete recipe they don't own", zap.String("recipe_id", id), zap.String("user_id", userID))
		return domain.ErrRecipeNotFound
	}

	_, err = r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"id":         &types.AttributeValueMemberS{Value: id},
			"created_at": &types.AttributeValueMemberS{Value: recipe.CreatedAt.Format(time.RFC3339Nano)},
		},
	})

	if err != nil {
		logger.Error(ctx, "Failed to delete recipe from DynamoDB", err, zap.String("recipe_id", id))
		return fmt.Errorf("failed to delete recipe: %w", err)
	}

	logger.Info(ctx, "Recipe deleted successfully", zap.String("recipe_id", id), zap.String("user_id", userID))
	return nil
}
