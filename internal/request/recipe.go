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
