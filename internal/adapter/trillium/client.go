package trillium

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// TrilliumClient handles communication with Trillium ETAPI.
type TrilliumClient struct {
	url        string
	token      string
	httpClient *http.Client
}

func NewClient(url, token string) *TrilliumClient {
	return &TrilliumClient{
		url:   url,
		token: token,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type createNotePayload struct {
	ParentNoteID string `json:"parentNoteId"`
	Title        string `json:"title"`
	Type         string `json:"type"`
	Content      string `json:"content"`
}

// CreateNote creates a new note in Trillium.
func (c *TrilliumClient) CreateNote(ctx context.Context, title, contentHTML string) error {
	// TODO: Support parentNoteId config? Defaulting to "root".
	// Trillium root note often has ID "root".
	payload := createNotePayload{
		ParentNoteID: "root",
		Title:        title,
		Type:         "text", // 'text' type is HTML content
		Content:      contentHTML,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal trillium payload: %w", err)
	}

	apiURL := fmt.Sprintf("%s/etapi/create-note", c.url)
	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewReader(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("trillium request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("trillium returned status: %s", resp.Status)
	}

	return nil
}
