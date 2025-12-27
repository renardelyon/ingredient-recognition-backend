package model

// RecipeRecommendation represents recipe suggestions based on ingredients
type RecipeRecommendation struct {
	Recipes            []Recipe `json:"recipes"`
	TotalRecipes       int      `json:"total_recipes"`
	GeneratedAt        string   `json:"generated_at"`
	IngredientCount    int      `json:"ingredient_count"`
	UsedIngredients    []string `json:"used_ingredients"`
	MissingIngredients []string `json:"missing_ingredients,omitempty"`
}

// Recipe represents a single recipe recommendation
type Recipe struct {
	Name         string   `json:"name"`
	Cuisine      string   `json:"cuisine"`
	CookingTime  string   `json:"cooking_time"`
	Difficulty   string   `json:"difficulty"`
	Ingredients  []string `json:"ingredients"`
	Instructions []string `json:"instructions"`
	Nutrition    string   `json:"nutrition,omitempty"`
	Tips         string   `json:"tips,omitempty"`
}
