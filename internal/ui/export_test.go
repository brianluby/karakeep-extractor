package ui

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/brianluby/karakeep-extractor/internal/core/domain"
)

func TestJSONExporter_Export(t *testing.T) {
	exporter := NewJSONExporter()
	var buf bytes.Buffer

	repos := []domain.ExtractedRepo{
		{RepoID: "test/repo", URL: "http://github.com/test/repo"},
	}

	if err := exporter.Export(repos, &buf); err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	var decoded []domain.ExtractedRepo
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("Failed to decode JSON output: %v", err)
	}

	if len(decoded) != 1 || decoded[0].RepoID != "test/repo" {
		t.Errorf("Unexpected JSON content: %v", decoded)
	}
}

func TestCSVExporter_Export(t *testing.T) {
	exporter := NewCSVExporter()
	var buf bytes.Buffer

	stars := 10
	repos := []domain.ExtractedRepo{
		{RepoID: "test/repo", URL: "http://github.com/test/repo", Stars: &stars},
	}

	if err := exporter.Export(repos, &buf); err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	output := buf.String()
	// Check Header
	if !strings.Contains(output, "Rank,RepoID,URL,Stars") {
		t.Error("Header missing or incorrect")
	}
	// Check Data
	if !strings.Contains(output, "1,test/repo,http://github.com/test/repo,10") {
		t.Error("Data row missing or incorrect")
	}
}

func TestGetExporter(t *testing.T) {
	e, err := GetExporter("json")
	if err != nil {
		t.Errorf("Expected no error for json, got %v", err)
	}
	if _, ok := e.(*JSONExporter); !ok {
		t.Error("Expected JSONExporter")
	}

	e, err = GetExporter("csv")
	if err != nil {
		t.Errorf("Expected no error for csv, got %v", err)
	}
	if _, ok := e.(*CSVExporter); !ok {
		t.Error("Expected CSVExporter")
	}

	e, err = GetExporter("table")
	if err != nil || e != nil {
		t.Errorf("Expected nil exporter and nil error for table, got %v, %v", e, err)
	}

	_, err = GetExporter("invalid")
	if err == nil {
		t.Error("Expected error for invalid format")
	}
}
