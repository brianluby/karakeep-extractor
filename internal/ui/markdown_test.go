package ui

import (
	"strings"
	"testing"
	"time"

	"github.com/brianluby/karakeep-extractor/internal/core/domain"
)

func TestMarkdownFormatter_FormatTable(t *testing.T) {
	formatter := NewMarkdownFormatter()
	now := time.Date(2023, 10, 27, 10, 0, 0, 0, time.UTC)
	stars := 100
	desc := "Test | Description"

	repos := []domain.ExtractedRepo{
		{
			RepoID:       "owner/repo",
			URL:          "http://github.com/owner/repo",
			Stars:        &stars,
			LastPushedAt: &now,
			Description:  &desc,
		},
	}

	output := formatter.FormatTable(repos)

	if !strings.Contains(output, "| Rank | Repository |") {
		t.Error("Header missing")
	}
	if !strings.Contains(output, "[owner/repo](http://github.com/owner/repo)") {
		t.Error("Link formatting incorrect")
	}
	if !strings.Contains(output, "Test \\| Description") {
		t.Errorf("Description pipe escaping incorrect. Output: %q", output)
	}
	if !strings.Contains(output, "2023-10-27") {
		t.Error("Date formatting incorrect")
	}
}
