package analysis

import (
	"strings"
	"testing"

	"github.com/brianluby/karakeep-extractor/internal/core/domain"
)

func TestBuildMessages(t *testing.T) {
	desc := "A test repo"
	lang := "Go"
	stars := 10
	forks := 2
	
	repos := []domain.ExtractedRepo{
		{
			RepoID:      "owner/repo",
			Description: &desc,
			Language:    &lang,
			Stars:       &stars,
			Forks:       &forks,
		},
	}

	msgs, err := BuildMessages("Summarize this", repos)
	if err != nil {
		t.Fatalf("BuildMessages failed: %v", err)
	}

	if len(msgs) != 2 {
		t.Fatalf("Expected 2 messages, got %d", len(msgs))
	}

	if msgs[0].Role != "system" {
		t.Errorf("Expected system role, got %s", msgs[0].Role)
	}
	if !strings.Contains(msgs[0].Content, "owner/repo") {
		t.Errorf("Expected system message to contain repo name")
	}
	if !strings.Contains(msgs[0].Content, "Go") {
		t.Errorf("Expected system message to contain language")
	}

	if msgs[1].Role != "user" {
		t.Errorf("Expected user role, got %s", msgs[1].Role)
	}
	if msgs[1].Content != "Summarize this" {
		t.Errorf("Expected user message content to match query")
	}
}
