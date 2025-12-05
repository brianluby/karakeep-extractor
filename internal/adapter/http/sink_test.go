package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brianluby/karakeep-extractor/internal/core/domain"
)

func TestHTTPSink_Send(t *testing.T) {
	// Mock Server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify Method
		if r.Method != "POST" {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		// Verify Headers
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type application/json")
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("Expected Authorization header")
		}

		// Verify Body
		var repos []domain.ExtractedRepo
		if err := json.NewDecoder(r.Body).Decode(&repos); err != nil {
			t.Errorf("Failed to decode body: %v", err)
		}
		if len(repos) != 1 || repos[0].RepoID != "test/repo" {
			t.Errorf("Unexpected payload: %v", repos)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Test Sink
	headers := []string{"Authorization: Bearer test-token"}
	sink := NewHTTPSink(server.URL, headers)

	repos := []domain.ExtractedRepo{
		{RepoID: "test/repo", URL: "http://github.com/test/repo"},
	}

	if err := sink.Send(context.Background(), repos); err != nil {
		t.Fatalf("Send failed: %v", err)
	}
}

func TestHTTPSink_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	sink := NewHTTPSink(server.URL, nil)
	if err := sink.Send(context.Background(), nil); err == nil {
		t.Error("Expected error on 500, got nil")
	}
}
