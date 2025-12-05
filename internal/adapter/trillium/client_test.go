package trillium

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTrilliumClient_CreateNote(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/etapi/create-note" {
			t.Errorf("Expected path /etapi/create-note, got %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		if r.Header.Get("Authorization") != "test-token" {
			t.Errorf("Expected Authorization header")
		}

		var payload createNotePayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Errorf("Failed to decode payload: %v", err)
		}
		if payload.Title != "Test Title" {
			t.Errorf("Expected title 'Test Title', got %s", payload.Title)
		}
		if payload.Content != "<b>Hello</b>" {
			t.Errorf("Expected content '<b>Hello</b>', got %s", payload.Content)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	err := client.CreateNote(context.Background(), "Test Title", "<b>Hello</b>")
	if err != nil {
		t.Fatalf("CreateNote failed: %v", err)
	}
}
