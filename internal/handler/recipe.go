package handler

import (
	"net/http"

	"ingredient-recognition-backend/internal/middleware"
	"ingredient-recognition-backend/internal/request"
	"ingredient-recognition-backend/internal/service"
	"ingredient-recognition-backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RecipeHandler struct {
	recipeService service.RecipeService
}

func NewRecipeHandler(recipeService service.RecipeService) *RecipeHandler {
	return &RecipeHandler{
		recipeService: recipeService,
	}
}

// RecommendRecipes generates recipe recommendations based on provided ingredients
// POST /api/recipes/recommend
func (h *RecipeHandler) RecommendRecipes(c *gin.Context) {
	logger.Info(c.Request.Context(), "Recipe recommendation request received")

	var req request.RecommendRecipesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn(c.Request.Context(), "Invalid recipe recommendation request", zap.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ingredients are required"})
		return
	}

	logger.Debug(c.Request.Context(), "Processing recipe recommendation", zap.Int("ingredient_count", len(req.Ingredients)))

	recommendation, err := h.recipeService.RecommendRecipes(c.Request.Context(), req.Ingredients)
	if err != nil {
		logger.Error(c.Request.Context(), "Recipe recommendation service failed", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate recipes"})
		return
	}

	logger.Info(c.Request.Context(), "Recipe recommendation completed", zap.Int("recipe_count", len(recommendation.Recipes)))
	c.JSON(http.StatusOK, recommendation)
}

// SaveRecipe saves a recipe for the authenticated user
// POST /api/v1/recipes/saved
func (h *RecipeHandler) SaveRecipe(c *gin.Context) {
	logger.Info(c.Request.Context(), "Save recipe request received")

	// Get user ID from context (set by auth middleware)
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		logger.Warn(c.Request.Context(), "Failed to get user from context", zap.String("error", err.Error()))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req request.SaveRecipeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn(c.Request.Context(), "Invalid save recipe request", zap.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	logger.Debug(c.Request.Context(), "Processing save recipe request",
		zap.String("user_id", userID),
		zap.String("recipe_name", req.Name))

	recipe, err := h.recipeService.SaveRecipe(c.Request.Context(), userID, &req)
	if err != nil {
		logger.Error(c.Request.Context(), "Failed to save recipe", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save recipe"})
		return
	}

	logger.Info(c.Request.Context(), "Recipe saved successfully",
		zap.String("recipe_id", recipe.ID),
		zap.String("user_id", userID))
	c.JSON(http.StatusCreated, recipe)
}

// GetUserRecipes retrieves all saved recipes for the authenticated user
// GET /api/v1/recipes/saved
func (h *RecipeHandler) GetUserRecipes(c *gin.Context) {
	logger.Info(c.Request.Context(), "Get user recipes request received")

	// Get user ID from context (set by auth middleware)
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		logger.Warn(c.Request.Context(), "Failed to get user from context", zap.String("error", err.Error()))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	recipes, err := h.recipeService.GetUserRecipes(c.Request.Context(), userID)
	if err != nil {
		logger.Error(c.Request.Context(), "Failed to get user recipes", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get recipes"})
		return
	}

	logger.Info(c.Request.Context(), "Retrieved user recipes",
		zap.String("user_id", userID),
		zap.Int("count", len(recipes)))
	c.JSON(http.StatusOK, gin.H{"recipes": recipes, "total": len(recipes)})
}

// GetRecipeByID retrieves a specific saved recipe by ID
// GET /api/v1/recipes/saved/:id
func (h *RecipeHandler) GetRecipeByID(c *gin.Context) {
	logger.Info(c.Request.Context(), "Get recipe by ID request received")

	// Get user ID from context (set by auth middleware)
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		logger.Warn(c.Request.Context(), "Failed to get user from context", zap.String("error", err.Error()))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	recipeID := c.Param("id")
	if recipeID == "" {
		logger.Warn(c.Request.Context(), "Recipe ID not provided")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Recipe ID is required"})
		return
	}

	recipe, err := h.recipeService.GetRecipeByID(c.Request.Context(), recipeID, userID)
	if err != nil {
		logger.Error(c.Request.Context(), "Failed to get recipe", err, zap.String("recipe_id", recipeID))
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
		return
	}

	logger.Info(c.Request.Context(), "Retrieved recipe",
		zap.String("recipe_id", recipeID),
		zap.String("user_id", userID))
	c.JSON(http.StatusOK, recipe)
}

// DeleteRecipe deletes a saved recipe
// DELETE /api/v1/recipes/saved/:id
func (h *RecipeHandler) DeleteRecipe(c *gin.Context) {
	logger.Info(c.Request.Context(), "Delete recipe request received")

	// Get user ID from context (set by auth middleware)
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		logger.Warn(c.Request.Context(), "Failed to get user from context", zap.String("error", err.Error()))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	recipeID := c.Param("id")
	if recipeID == "" {
		logger.Warn(c.Request.Context(), "Recipe ID not provided")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Recipe ID is required"})
		return
	}

	err = h.recipeService.DeleteRecipe(c.Request.Context(), recipeID, userID)
	if err != nil {
		logger.Error(c.Request.Context(), "Failed to delete recipe", err, zap.String("recipe_id", recipeID))
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
		return
	}

	logger.Info(c.Request.Context(), "Recipe deleted successfully",
		zap.String("recipe_id", recipeID),
		zap.String("user_id", userID))
	c.JSON(http.StatusOK, gin.H{"message": "Recipe deleted successfully"})
}
