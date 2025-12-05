package ui

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/brianluby/karakeep-extractor/internal/core/domain"
)

// JSONExporter exports repositories as a JSON array.
type JSONExporter struct{}

func NewJSONExporter() *JSONExporter {
	return &JSONExporter{}
}

func (j *JSONExporter) Export(repos []domain.ExtractedRepo, w io.Writer) error {
	// Use indentation for better readability, or standard? Spec says "valid JSON".
	// Indentation is friendlier for CLI usage.
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(repos)
}

// CSVExporter exports repositories as a CSV file with headers.
type CSVExporter struct{}

func NewCSVExporter() *CSVExporter {
	return &CSVExporter{}
}

func (c *CSVExporter) Export(repos []domain.ExtractedRepo, w io.Writer) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	// Write Header
	header := []string{"Rank", "RepoID", "URL", "Stars", "Forks", "LastPushedAt", "Description", "Language"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	for i, repo := range repos {
		rank := strconv.Itoa(i + 1)
		stars := "0"
		if repo.Stars != nil {
			stars = strconv.Itoa(*repo.Stars)
		}
		forks := "0"
		if repo.Forks != nil {
			forks = strconv.Itoa(*repo.Forks)
		}
		lastPushed := ""
		if repo.LastPushedAt != nil {
			lastPushed = repo.LastPushedAt.Format("2006-01-02T15:04:05Z")
		}
		desc := ""
		if repo.Description != nil {
			desc = *repo.Description
		}
		lang := ""
		if repo.Language != nil {
			lang = *repo.Language
		}

		record := []string{
			rank,
			repo.RepoID,
			repo.URL,
			stars,
			forks,
			lastPushed,
			desc,
			lang,
		}

		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write CSV record for %s: %w", repo.RepoID, err)
		}
	}

	return nil
}

// GetExporter returns the appropriate exporter based on the format string.
// Returns nil if format is "table" (handled by TableRenderer directly for now) or unknown.
func GetExporter(format string) (domain.Exporter, error) {
	switch format {
	case "json":
		return NewJSONExporter(), nil
	case "csv":
		return NewCSVExporter(), nil
	case "table":
		return nil, nil // Special case, existing logic
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}
