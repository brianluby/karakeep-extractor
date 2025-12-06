package ui

import (
	"fmt"
	"strings"
	"html"

	"github.com/brianluby/karakeep-extractor/internal/core/domain"
)

// MarkdownFormatter formats repositories into a Markdown table string.
type MarkdownFormatter struct{}

func NewMarkdownFormatter() *MarkdownFormatter {
	return &MarkdownFormatter{}
}

func (m *MarkdownFormatter) FormatTable(repos []domain.ExtractedRepo) string {
	var b strings.Builder

	// Header
	b.WriteString("| Rank | Repository | Stars | Forks | Last Updated | Description |\n")
	b.WriteString("|------|------------|-------|-------|--------------|-------------|\n")

	for i, repo := range repos {
		rank := i + 1
		stars := 0
		if repo.Stars != nil {
			stars = *repo.Stars
		}
		forks := 0
		if repo.Forks != nil {
			forks = *repo.Forks
		}
		updated := "-"
		if repo.LastPushedAt != nil {
			updated = repo.LastPushedAt.Format("2006-01-02")
		}
		
		desc := ""
		if repo.Description != nil {
			desc = html.EscapeString(*repo.Description)
			// Escape pipes in description to avoid breaking table
			desc = strings.ReplaceAll(desc, "|", "\\|")
			// Truncate if too long? Or keep full. Trillium handles scroll.
		}

		// Link the repo name
		nameLink := fmt.Sprintf("[%s](%s)", repo.RepoID, repo.URL)

		fmt.Fprintf(&b, "| %d | %s | %d | %d | %s | %s |\n", 
			rank, nameLink, stars, forks, updated, desc)
	}

	return b.String()
}
