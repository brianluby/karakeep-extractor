package domain

// LLMConfig holds the configuration for the LLM provider.
type LLMConfig struct {
	Provider  string `yaml:"provider"`
	BaseURL   string `yaml:"base_url"`
	APIKey    string `yaml:"api_key"`
	Model     string `yaml:"model"`
	MaxTokens int    `yaml:"max_tokens,omitempty"`
}

// RepositoryContext represents a subset of repository data for LLM analysis.
type RepositoryContext struct {
	Name        string   `json:"name"`
	URL         string   `json:"url"`
	Description string   `json:"description"`
	Language    string   `json:"language"`
	Stars       int      `json:"stars"`
	Forks       int      `json:"forks"`
	LastUpdated string   `json:"last_updated,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

// Message represents a single message in the chat completion conversation.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// AnalysisRequest represents the payload sent to the LLM API.
type AnalysisRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	// OpenAI specific, but generic enough
	MaxTokens int `json:"max_tokens,omitempty"`
}
