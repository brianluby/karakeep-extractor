package http

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

type HTTPSink struct {
	url     string
	headers map[string]string
	client  *http.Client
}

func NewHTTPSink(url string, headers []string) *HTTPSink {
	return &HTTPSink{
		url:     url,
		headers: parseHeaders(headers),
		client: &http.Client{
			Timeout: 30 * time.Second, // Generous timeout for large payloads
		},
	}
}

func (s *HTTPSink) Send(ctx context.Context, repos []domain.ExtractedRepo) error {
	payload, err := json.Marshal(repos)
	if err != nil {
		return fmt.Errorf("failed to marshal repos for sink: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.url, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("failed to create sink request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	for k, v := range s.headers {
		req.Header.Set(k, v)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("sink request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("sink returned status: %s", resp.Status)
	}

	return nil
}

func parseHeaders(headers []string) map[string]string {
	m := make(map[string]string)
	for _, h := range headers {
		parts := strings.SplitN(h, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			if key != "" {
				m[key] = val
			}
		}
	}
	return m
}