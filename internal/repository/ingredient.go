package repository

import (
	"fmt"
	"ingredient-recognition-backend/internal/domain"
)

// IngredientRepository defines the methods for interacting with the ingredient data source.
type IngredientRepository interface {
	Save(ingredient domain.Ingredient) error
	FindByID(id string) (domain.Ingredient, error)
	FindAll() ([]domain.Ingredient, error)
}

// InMemoryIngredientRepository is a concrete implementation of IngredientRepository that stores ingredients in memory.
type InMemoryIngredientRepository struct {
	ingredients map[string]domain.Ingredient
}

// NewInMemoryIngredientRepository creates a new instance of InMemoryIngredientRepository.
func NewInMemoryIngredientRepository() *InMemoryIngredientRepository {
	return &InMemoryIngredientRepository{
		ingredients: make(map[string]domain.Ingredient),
	}
}

// Save saves an ingredient to the in-memory store.
func (repo *InMemoryIngredientRepository) Save(ingredient domain.Ingredient) error {
	repo.ingredients[ingredient.ID] = ingredient
	return nil
}

// FindByID retrieves an ingredient by its ID from the in-memory store.
func (repo *InMemoryIngredientRepository) FindByID(id string) (domain.Ingredient, error) {
	ingredient, exists := repo.ingredients[id]
	if !exists {
		return domain.Ingredient{}, fmt.Errorf("ingredient with ID %s not found", id)
	}
	return ingredient, nil
}

// FindAll retrieves all ingredients from the in-memory store.
func (repo *InMemoryIngredientRepository) FindAll() ([]domain.Ingredient, error) {
	var allIngredients []domain.Ingredient
	for _, ingredient := range repo.ingredients {
		allIngredients = append(allIngredients, ingredient)
	}
	return allIngredients, nil
}
