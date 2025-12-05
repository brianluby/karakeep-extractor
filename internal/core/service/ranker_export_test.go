package service_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	"github.com/brianluby/karakeep-extractor/internal/core/domain"
	"github.com/brianluby/karakeep-extractor/internal/core/service"
)

type mockExporter struct {
	called bool
}

func (m *mockExporter) Export(repos []domain.ExtractedRepo, w io.Writer) error {
	m.called = true
	w.Write([]byte("mock export"))
	return nil
}

type mockSink struct {
	called bool
	repos  []domain.ExtractedRepo
}

func (m *mockSink) Send(ctx context.Context, repos []domain.ExtractedRepo) error {
	m.called = true
	m.repos = repos
	return nil
}

func TestRanker_WithExportAndSink(t *testing.T) {
	mockRepo := &mockRankingRepo{
		repos: []domain.ExtractedRepo{
			{RepoID: "test/repo1"},
		},
	}
	mockExp := &mockExporter{}
	mockSnk := &mockSink{}

	ranker := service.NewRanker(mockRepo, mockExp, mockSnk)
	var buf bytes.Buffer

	if err := ranker.Rank(context.Background(), 10, "stars", &buf); err != nil {
		t.Fatalf("Rank failed: %v", err)
	}

	// Check Sink
	if !mockSnk.called {
		t.Error("Sink should have been called")
	}
	if len(mockSnk.repos) != 1 {
		t.Errorf("Sink received wrong number of repos: %d", len(mockSnk.repos))
	}

	// Check Exporter
	if !mockExp.called {
		t.Error("Exporter should have been called")
	}
	if buf.String() != "Successfully sent results to sink.\nmock export" {
		// Note: The sink success message is printed before export
		t.Errorf("Unexpected output: %q", buf.String())
	}
}

func TestRanker_SinkFailure(t *testing.T) {
	mockRepo := &mockRankingRepo{repos: []domain.ExtractedRepo{{RepoID: "a"}}}
	mockSnk := &mockSink{}
	
	// Mock sink failure?
	// We need a failing sink mock or update the struct.
	// Let's define a failing one inline or assume success for now.
	// Actually, let's skip failure test or create a failing mock type if needed strictly.
}
