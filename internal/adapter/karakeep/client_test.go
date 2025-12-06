package karakeep_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/brianluby/karakeep-extractor/internal/adapter/karakeep"
	"github.com/brianluby/karakeep-extractor/internal/core/domain"
)

func TestNewClientAndAuth(t *testing.T) {
	// Mock server for successful connection
	serverSuccess := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{"status":"ok"}`)
	}))
	defer serverSuccess.Close()

	// Mock server for unauthorized
	serverUnauthorized := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer serverUnauthorized.Close()

	// Test case 1: Successful connection
	cfg := &domain.KarakeepConfig{
		BaseURL:  serverSuccess.URL,
		APIToken: "test-token",
	}
	client := karakeep.NewClient(cfg)
	resp, err := client.HTTPClient.Get(serverSuccess.URL)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %d", resp.StatusCode)
	}

	// Test case 2: Unauthorized
	cfgUnauthorized := &domain.KarakeepConfig{
		BaseURL:  serverUnauthorized.URL,
		APIToken: "wrong-token",
	}
	clientUnauthorized := karakeep.NewClient(cfgUnauthorized)
	respUnauthorized, errUnauthorized := clientUnauthorized.HTTPClient.Get(serverUnauthorized.URL)
	if errUnauthorized != nil {
		t.Fatalf("Expected no error, got %v", errUnauthorized)
	}
	if respUnauthorized.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status Unauthorized, got %d", respUnauthorized.StatusCode)
	}
}

func TestRetryLogic(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 { // Fail first 2 attempts
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{"status":"ok"}`)
	}))
	defer server.Close()

	cfg := &domain.KarakeepConfig{
		BaseURL:  server.URL,
		APIToken: "test-token",
	}
	client := karakeep.NewClient(cfg) // NewClient needs to wrap HTTPClient with retry logic

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, "GET", server.URL, nil)
	resp, err := client.HTTPClient.Do(req) // HTTPClient should have retry logic applied
	if err != nil {
		t.Fatalf("Expected no error after retries, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK after retries, got %d", resp.StatusCode)
	}
	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
}

func TestFetchBookmarks(t *testing.T) {
	// Mock server for successful bookmark fetch
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// Return all bookmarks in one go, wrapped in "bookmarks"
		response := struct {
			Bookmarks []domain.RawBookmark `json:"bookmarks"`
		}{
			Bookmarks: []domain.RawBookmark{
				{
					ID: "1",
					Content: struct {
						URL         string `json:"url"`
						Title       string `json:"title"`
						Description string `json:"description"`
					}{URL: "https://github.com/repo1", Title: "Repo 1", Description: ""},
				},
				{
					ID: "2",
					Content: struct {
						URL         string `json:"url"`
						Title       string `json:"title"`
						Description string `json:"description"`
					}{URL: "https://example.com/article1", Title: "Article 1", Description: ""},
				},
				{
					ID: "3",
					Content: struct {
						URL         string `json:"url"`
						Title       string `json:"title"`
						Description string `json:"description"`
					}{URL: "https://github.com/repo2", Title: "Repo 2", Description: ""},
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	cfg := &domain.KarakeepConfig{
		BaseURL:  server.URL,
		APIToken: "test-token",
	}
	client := karakeep.NewClient(cfg)

	// Test case: Fetch all bookmarks
	bookmarks, err := client.FetchBookmarks(context.Background())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(bookmarks) != 3 {
		t.Errorf("Expected 3 bookmarks, got %d", len(bookmarks))
	}
	if bookmarks[0].Content.URL != "https://github.com/repo1" {
		t.Errorf("Expected URL https://github.com/repo1, got %s", bookmarks[0].Content.URL)
	}
	if bookmarks[2].Content.URL != "https://github.com/repo2" {
		t.Errorf("Expected URL https://github.com/repo2, got %s", bookmarks[2].Content.URL)
	}
}
