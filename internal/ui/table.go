package ui

import (
	"fmt"
	"io"
	"text/tabwriter"
	"time"

	"github.com/brianluby/karakeep-extractor/internal/core/domain"
)

// TableRenderer renders a list of repos as a formatted table.
type TableRenderer struct {
	writer *tabwriter.Writer
}

func NewTableRenderer(output io.Writer) *TableRenderer {
	// MinWidth=0, TabWidth=8, Padding=2, PadChar=' ', Flags=0
	return &TableRenderer{
		writer: tabwriter.NewWriter(output, 0, 8, 2, ' ', 0),
	}
}

// Render prints the table to the configured writer.
func (t *TableRenderer) Render(repos []domain.ExtractedRepo) error {
	// Header
	fmt.Fprintln(t.writer, "RANK\tNAME\tSTARS\tFORKS\tUPDATED")

	for i, repo := range repos {
		rank := i + 1
		name := repo.RepoID
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
			updated = formatRelativeTime(*repo.LastPushedAt)
		}

		fmt.Fprintf(t.writer, "%d\t%s\t%d\t%d\t%s\n", rank, name, stars, forks, updated)
	}

	return t.writer.Flush()
}

func formatRelativeTime(t time.Time) string {
	diff := time.Since(t)
	
	switch {
	case diff < time.Hour:
		return fmt.Sprintf("%dm ago", int(diff.Minutes()))
	case diff < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(diff.Hours()))
	case diff < 30*24*time.Hour:
		return fmt.Sprintf("%dd ago", int(diff.Hours()/24))
	default:
		// For older dates, maybe months or just YYYY-MM-DD
		return t.Format("2006-01-02")
	}
}
