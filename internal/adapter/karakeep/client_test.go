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
		page := r.URL.Query().Get("page")
		if page == "" || page == "1" {
			json.NewEncoder(w).Encode([]domain.RawBookmark{
				{ID: "1", URL: "https://github.com/repo1", Title: "Repo 1", Content: ""},
				{ID: "2", URL: "https://example.com/article1", Title: "Article 1", Content: ""},
			})
		} else if page == "2" {
			json.NewEncoder(w).Encode([]domain.RawBookmark{
				{ID: "3", URL: "https://github.com/repo2", Title: "Repo 2", Content: ""},
			})
		} else {
			json.NewEncoder(w).Encode([]domain.RawBookmark{})
		}
	}))
	defer server.Close()

	cfg := &domain.KarakeepConfig{
		BaseURL:  server.URL,
		APIToken: "test-token",
	}
	client := karakeep.NewClient(cfg)

	// Test case 1: Fetch page 1
	bookmarks, err := client.FetchBookmarks(context.Background(), 1)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(bookmarks) != 2 {
		t.Errorf("Expected 2 bookmarks, got %d", len(bookmarks))
	}
	if bookmarks[0].URL != "https://github.com/repo1" {
		t.Errorf("Expected URL https://github.com/repo1, got %s", bookmarks[0].URL)
	}

	// Test case 2: Fetch page 2
	bookmarks2, err2 := client.FetchBookmarks(context.Background(), 2)
	if err2 != nil {
		t.Fatalf("Expected no error, got %v", err2)
	}
	if len(bookmarks2) != 1 {
		t.Errorf("Expected 1 bookmark, got %d", len(bookmarks2))
	}
	if bookmarks2[0].URL != "https://github.com/repo2" {
		t.Errorf("Expected URL https://github.com/repo2, got %s", bookmarks2[0].URL)
	}

	// Test case 3: Fetch empty page
	bookmarks3, err3 := client.FetchBookmarks(context.Background(), 3)
	if err3 != nil {
		t.Fatalf("Expected no error, got %v", err3)
	}
	if len(bookmarks3) != 0 {
		t.Errorf("Expected 0 bookmarks, got %d", len(bookmarks3))
	}
}
