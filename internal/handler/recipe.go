package handler

import (
	"net/http"

	"ingredient-recognition-backend/internal/request"
	"ingredient-recognition-backend/internal/service"

	"github.com/gin-gonic/gin"
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
	var req request.RecommendRecipesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ingredients are required"})
		return
	}

	recommendation, err := h.recipeService.RecommendRecipes(c.Request.Context(), req.Ingredients)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate recipes"})
		return
	}

	c.JSON(http.StatusOK, recommendation)
}
