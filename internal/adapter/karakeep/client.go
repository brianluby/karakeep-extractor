package karakeep

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/brianluby/karakeep-extractor/internal/core/domain"
)

const (
	maxRetries    = 3
	initialBackoff = 100 * time.Millisecond
)

type Client struct {
	Config     *domain.KarakeepConfig
	HTTPClient *http.Client
}

// NewClient creates a new Karakeep API client with retry logic.
func NewClient(cfg *domain.KarakeepConfig) *Client {
	// Custom HTTP client with retry logic
	httpClient := &http.Client{
		Transport: &authTransport{
			Transport: &retryTransport{
				roundTripper: http.DefaultTransport,
				maxRetries:   maxRetries,
				initialDelay: initialBackoff,
			},
			token: cfg.APIToken,
		},
		Timeout: 10 * time.Second, // Global timeout for Karakeep API calls
	}

	return &Client{
		Config:     cfg,
		HTTPClient: httpClient,
	}
}

// authTransport adds the Authorization header to requests.
type authTransport struct {
	Transport http.RoundTripper
	token     string
}

func (at *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+at.token)
	return at.Transport.RoundTrip(req)
}


// retryTransport implements http.RoundTripper with exponential backoff.
type retryTransport struct {
	roundTripper http.RoundTripper
	maxRetries   int
	initialDelay time.Duration
}

func (rt *retryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var (
		resp *http.Response
		err  error
		delay time.Duration
	)

	for i := 0; i < rt.maxRetries; i++ {
		// Clone the request for each retry to avoid issues with body being read
		reqClone := req.Clone(req.Context())
		if req.GetBody != nil { // If body is available, reset it for the clone
			bodyReader, bodyErr := req.GetBody()
			if bodyErr != nil {
				return nil, bodyErr
			}
			reqClone.Body = bodyReader
		}

		resp, err = rt.roundTripper.RoundTrip(reqClone)
		if err == nil {
			if resp.StatusCode != http.StatusTooManyRequests {
				return resp, nil // Success or non-retryable error
			}
		} else {
			// Network error, context cancelled, etc.
			// Decide if these should be retried. For now, only 429 is explicitly retried.
			// Other errors are returned immediately after first attempt.
			return resp, err 
		}

		// Handle retryable errors (StatusTooManyRequests)
		if resp.StatusCode == http.StatusTooManyRequests {
			resp.Body.Close() // Close body to reuse connection
			delay = rt.initialDelay * time.Duration(1<<uint(i)) // Exponential backoff
			select {
			case <-req.Context().Done():
				return nil, req.Context().Err()
			case <-time.After(delay):
				// Continue to next retry attempt
			}
			continue
		}
	}

	// If all retries fail
	if resp != nil {
		return resp, fmt.Errorf("failed after %d retries, last status: %d", rt.maxRetries, resp.StatusCode)
	}
	// This case would be for network errors that were not retried or if initial attempt also failed with a non-429
	return nil, fmt.Errorf("request failed after %d retries (check connection/server status), last error: %v", rt.maxRetries, err)
}

// Ping verifies connectivity and authentication to the Karakeep API.
func (c *Client) Ping(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", c.Config.BaseURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return handleErrorResponse(resp)
	}
	return nil
}

// handleErrorResponse maps HTTP status codes to user-friendly errors.
func handleErrorResponse(resp *http.Response) error {
	switch resp.StatusCode {
	case http.StatusUnauthorized:
		return fmt.Errorf("authentication failed: invalid Karakeep token")
	case http.StatusNotFound:
		return fmt.Errorf("API endpoint not found: %s", resp.Request.URL.String())
	case http.StatusInternalServerError:
		return fmt.Errorf("Karakeep server error: %s", resp.Status)
	case http.StatusBadRequest:
		return fmt.Errorf("bad request to Karakeep API: %s", resp.Status)
	default:
		return fmt.Errorf("Karakeep API returned unexpected status: %s", resp.Status)
	}
}


// FetchBookmarks fetches bookmarks from the Karakeep API.
func (c *Client) FetchBookmarks(ctx context.Context) ([]domain.RawBookmark, error) {
	baseURL := strings.TrimSuffix(c.Config.BaseURL, "/")
	var allBookmarks []domain.RawBookmark
	var cursor string
	limit := 100

	for {
		url := fmt.Sprintf("%s/bookmarks?archived=false&includeContent=true&limit=%d", baseURL, limit)
		if cursor != "" {
			url += fmt.Sprintf("&cursor=%s", cursor)
		}

		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch bookmarks: %w", err)
		}
		// Defer inside loop is okay if loop isn't huge, but better to close explicitly
		// We'll read properly below

		bodyBytes, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			return nil, handleErrorResponse(resp)
		}

		var response struct {
			Bookmarks  []domain.RawBookmark `json:"bookmarks"`
			NextCursor *string              `json:"nextCursor"`
		}

		if err := json.Unmarshal(bodyBytes, &response); err != nil {
			return nil, fmt.Errorf("failed to decode bookmarks: %w", err)
		}

		if len(response.Bookmarks) > 0 {
			allBookmarks = append(allBookmarks, response.Bookmarks...)
		}

		if response.NextCursor == nil || *response.NextCursor == "" {
			break
		}
		cursor = *response.NextCursor
	}

	return allBookmarks, nil
}
