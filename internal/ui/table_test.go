package ui

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/brianluby/karakeep-extractor/internal/core/domain"
)

func TestTableRenderer_Render(t *testing.T) {
	var buf bytes.Buffer
	renderer := NewTableRenderer(&buf)

	now := time.Now()
	repo1 := domain.ExtractedRepo{
		RepoID:       "owner/repo1",
		Stars:        intPtr(100),
		Forks:        intPtr(50),
		LastPushedAt: &now,
	}
	
	repos := []domain.ExtractedRepo{repo1}

	if err := renderer.Render(repos); err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "RANK") || !strings.Contains(output, "owner/repo1") {
		t.Errorf("Output missing expected headers or data: %s", output)
	}
	
	// Check relative time formatting
	if !strings.Contains(output, "0m ago") { // Roughly immediate
		t.Errorf("Expected relative time formatting, got output: %s", output)
	}
}

func intPtr(i int) *int {
	return &i
}
