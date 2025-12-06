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

	ingredients, err := h.detectorService.DetectIngredientsFromImage(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to detect ingredients"})
		return
	}

	c.JSON(http.StatusOK, ingredients)
}
