package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/brianluby/karakeep-extractor/internal/core/domain"
)

type Client struct {
	token      string
	baseURL    string
	httpClient *http.Client
}

func NewClient(token string) *Client {
	return &Client{
		token:   token,
		baseURL: "https://api.github.com",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// WithBaseURL allows overriding the base URL for testing.
func (c *Client) WithBaseURL(url string) *Client {
	c.baseURL = url
	return c
}

type githubRepoResponse struct {
	StargazersCount int       `json:"stargazers_count"`
	ForksCount      int       `json:"forks_count"`
	PushedAt        time.Time `json:"pushed_at"`
	Description     string    `json:"description"`
	Language        string    `json:"language"`
}

func (c *Client) GetRepoStats(ctx context.Context, owner, repo string) (*domain.RepoStats, int, error) {
	url := fmt.Sprintf("%s/repos/%s/%s", c.baseURL, owner, repo)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	if c.token != "" {
		req.Header.Set("Authorization", "token "+c.token)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Parse Rate Limit
	remaining := 0
	if remStr := resp.Header.Get("X-RateLimit-Remaining"); remStr != "" {
		remaining, _ = strconv.Atoi(remStr)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, remaining, domain.ErrRepoNotFound
	}
	if resp.StatusCode == http.StatusForbidden && remaining == 0 {
		return nil, 0, domain.ErrRateLimitExceeded
	}
	if resp.StatusCode != http.StatusOK {
		return nil, remaining, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var ghResp githubRepoResponse
	if err := json.NewDecoder(resp.Body).Decode(&ghResp); err != nil {
		return nil, remaining, fmt.Errorf("failed to decode response: %w", err)
	}

	stats := &domain.RepoStats{
		Stars:       ghResp.StargazersCount,
		Forks:       ghResp.ForksCount,
		LastPushed:  ghResp.PushedAt,
		Description: ghResp.Description,
		Language:    ghResp.Language,
	}

	return stats, remaining, nil
}
