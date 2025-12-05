package github

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_GetRepoStats(t *testing.T) {
	tests := []struct {
		name          string
		token         string
		handler       func(w http.ResponseWriter, r *http.Request)
		expectedStars int
		expectErr     bool
		expectRateLim int
	}{
		{
			name:  "Success",
			token: "valid-token",
			handler: func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("Authorization") != "token valid-token" {
					t.Errorf("Expected Authorization header")
				}
				w.Header().Set("X-RateLimit-Remaining", "4999")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{
					"stargazers_count": 100,
					"forks_count": 20,
					"pushed_at": "2023-01-01T12:00:00Z",
					"description": "Test Repo",
					"language": "Go"
				}`))
			},
			expectedStars: 100,
			expectRateLim: 4999,
		},
		{
			name: "Not Found",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-RateLimit-Remaining", "59")
				w.WriteHeader(http.StatusNotFound)
			},
			expectErr:     true,
			expectRateLim: 59,
		},
		{
			name: "Rate Limit Exceeded",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-RateLimit-Remaining", "0")
				w.WriteHeader(http.StatusForbidden)
			},
			expectErr:     true,
			expectRateLim: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(tt.handler))
			defer ts.Close()

			client := NewClient(tt.token).WithBaseURL(ts.URL)
			stats, rem, err := client.GetRepoStats(context.Background(), "owner", "repo")

			if tt.expectErr {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if stats.Stars != tt.expectedStars {
					t.Errorf("Expected stars %d, got %d", tt.expectedStars, stats.Stars)
				}
			}

			if rem != tt.expectRateLim {
				t.Errorf("Expected rate limit %d, got %d", tt.expectRateLim, rem)
			}
		})
	}
}