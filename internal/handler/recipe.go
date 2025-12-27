package handler

import (
	"net/http"

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
