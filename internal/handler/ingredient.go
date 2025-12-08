package handler

import (
	"net/http"

	"ingredient-recognition-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type IngredientHandler struct {
	detectorService service.DetectorService
}

func NewIngredientHandler(detectorService service.DetectorService) *IngredientHandler {
	return &IngredientHandler{detectorService: detectorService}
}

func (h *IngredientHandler) DetectIngredients(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image file"})
		return
	}

	ingredients, err := h.detectorService.DetectIngredientsFromImage(c.Request.Context(), file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to detect ingredients", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ingredients": ingredients})
}

// DetectIngredientsWithCustomLabels detects ingredients using a trained custom labels model
func (h *IngredientHandler) DetectIngredientsWithCustomLabels(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image file"})
		return
	}

	ingredients, err := h.detectorService.DetectIngredientsFromImageWithCustomLabels(c.Request.Context(), file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to detect ingredients with custom labels", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ingredients": ingredients})
}
