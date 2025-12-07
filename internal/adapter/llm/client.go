package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/brianluby/karakeep-extractor/internal/core/domain"
)

type Client struct {
	config domain.LLMConfig
	http   *http.Client
}

func NewClient(cfg domain.LLMConfig) *Client {
	return &Client{
		config: cfg,
		http: &http.Client{
			Timeout: 60 * time.Second, // Generous timeout for LLM thinking
		},
	}
}

type openAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func (c *Client) SendMessage(ctx context.Context, req domain.AnalysisRequest) (string, error) {
	// Normalize BaseURL
	baseURL := strings.TrimRight(c.config.BaseURL, "/")
	url := fmt.Sprintf("%s/chat/completions", baseURL)

	// Ensure request uses config model if not set
	if req.Model == "" {
		req.Model = c.config.Model
	}
	// Max tokens override
	if c.config.MaxTokens > 0 {
		req.MaxTokens = c.config.MaxTokens
	}

	body, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if c.config.APIKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.config.APIKey)
	}

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("network error calling LLM: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Specific handling for 401
		if resp.StatusCode == http.StatusUnauthorized {
			return "", fmt.Errorf("authentication failed: check your API key in 'karakeep config llm'")
		}

		// Try to read error from body
		var errResp openAIResponse
		json.NewDecoder(resp.Body).Decode(&errResp)
		msg := "unknown error"
		if errResp.Error != nil {
			msg = errResp.Error.Message
		}
		return "", fmt.Errorf("LLM API error (status %d): %s", resp.StatusCode, msg)
	}

	var response openAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("empty response from LLM")
	}

	return response.Choices[0].Message.Content, nil
}
