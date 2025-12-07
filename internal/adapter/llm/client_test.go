package llm

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brianluby/karakeep-extractor/internal/core/domain"
)

func TestClient_SendMessage(t *testing.T) {
	// Mock Server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify Request
		if r.Method != "POST" {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Errorf("Expected Bearer token")
		}
		
		var req domain.AnalysisRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request: %v", err)
		}
		if req.Model != "gpt-test" {
			t.Errorf("Expected model gpt-test, got %s", req.Model)
		}

		// Response
		resp := map[string]interface{}{
			"choices": []map[string]interface{}{
				{
					"message": map[string]interface{}{
						"content": "Test Response",
					},
				},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	cfg := domain.LLMConfig{
		Provider: "openai",
		BaseURL:  server.URL, // Use mock server URL
		APIKey:   "test-key",
		Model:    "gpt-test",
	}

	client := NewClient(cfg)
	
	req := domain.AnalysisRequest{
		Model: "gpt-test",
		Messages: []domain.Message{
			{Role: "user", Content: "Hello"},
		},
	}

	resp, err := client.SendMessage(context.Background(), req)
	if err != nil {
		t.Fatalf("SendMessage failed: %v", err)
	}

	if resp != "Test Response" {
		t.Errorf("Expected 'Test Response', got '%s'", resp)
	}
}
