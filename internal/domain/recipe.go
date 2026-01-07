package domain

import (
	"errors"
	"time"
)

// SavedRecipe represents a recipe saved by a user
type SavedRecipe struct {
	ID           string    `json:"id" dynamodbav:"id"`
	UserID       string    `json:"user_id" dynamodbav:"user_id"`
	Name         string    `json:"name" dynamodbav:"name"`
	Cuisine      string    `json:"cuisine" dynamodbav:"cuisine"`
	CookingTime  string    `json:"cooking_time" dynamodbav:"cooking_time"`
	Difficulty   string    `json:"difficulty" dynamodbav:"difficulty"`
	Ingredients  []string  `json:"ingredients" dynamodbav:"ingredients"`
	Instructions []string  `json:"instructions" dynamodbav:"instructions"`
	Nutrition    string    `json:"nutrition,omitempty" dynamodbav:"nutrition,omitempty"`
	Tips         string    `json:"tips,omitempty" dynamodbav:"tips,omitempty"`
	CreatedAt    time.Time `json:"created_at" dynamodbav:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" dynamodbav:"updated_at"`
}

var (
	ErrRecipeNotFound      = errors.New("recipe not found")
	ErrRecipeAlreadyExists = errors.New("recipe already exists")
)
