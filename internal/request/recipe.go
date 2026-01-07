package request

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type CacheControl struct {
	Type string `json:"type"`
}

type System struct {
	Type         string        `json:"type"`
	Text         string        `json:"text,omitempty"`
	CacheControl *CacheControl `json:"cache_control,omitempty"`
}

type BedrockModelConfig struct {
	AnthropicVersion string    `json:"anthropic_version"`
	MaxTokens        int       `json:"max_tokens"`
	Messages         []Message `json:"messages"`
	System           []System  `json:"system,omitempty"`
}

// RecommendRecipesRequest represents the request to get recipe recommendations
type RecommendRecipesRequest struct {
	Ingredients []string `json:"ingredients" binding:"required,min=1"`
}

// SaveRecipeRequest represents the request to save a recipe
type SaveRecipeRequest struct {
	Name         string   `json:"name" binding:"required"`
	Cuisine      string   `json:"cuisine" binding:"required"`
	CookingTime  string   `json:"cooking_time" binding:"required"`
	Difficulty   string   `json:"difficulty" binding:"required"`
	Ingredients  []string `json:"ingredients" binding:"required,min=1"`
	Instructions []string `json:"instructions" binding:"required,min=1"`
	Nutrition    string   `json:"nutrition,omitempty"`
	Tips         string   `json:"tips,omitempty"`
}
