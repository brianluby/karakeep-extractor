package service

import (
	"context"
	"errors"
	"testing"

	"github.com/brianluby/karakeep-extractor/internal/core/domain"
)

// Mock Repository
type MockRepo struct {
	repos map[string]*domain.ExtractedRepo
}

func (m *MockRepo) Save(ctx context.Context, repo domain.ExtractedRepo) error { return nil }
func (m *MockRepo) Exists(ctx context.Context, repoID string) (bool, error)   { return true, nil }
func (m *MockRepo) GetReposForEnrichment(ctx context.Context, limit int, force bool) ([]*domain.ExtractedRepo, error) {
	var res []*domain.ExtractedRepo
	for _, r := range m.repos {
		res = append(res, r)
	}
	return res, nil
}
func (m *MockRepo) UpdateRepoEnrichment(ctx context.Context, update domain.RepoEnrichmentUpdate) error {
	if r, ok := m.repos[update.RepoID]; ok {
		r.EnrichmentStatus = update.EnrichmentStatus
		if update.Stats != nil {
			r.Stars = &update.Stats.Stars
		}
	}
	return nil
}

// Mock Client
type MockClient struct {
	stats     map[string]*domain.RepoStats
	failRepo  string
	failError error
}

func (m *MockClient) GetRepoStats(ctx context.Context, owner, repo string) (*domain.RepoStats, int, error) {
	id := owner + "/" + repo
	if id == m.failRepo {
		return nil, 0, m.failError
	}
	if s, ok := m.stats[id]; ok {
		return s, 5000, nil
	}
	return nil, 5000, errors.New("not found")
}

// MockReporter for testing
type mockReporter struct{}

func (m *mockReporter) Start(total int, title string)   {}
func (m *mockReporter) Increment()                      {}
func (m *mockReporter) SetStatus(status string)         {}
func (m *mockReporter) Log(message string)              {}
func (m *mockReporter) Error(err error)                 {}
func (m *mockReporter) Finish(summary string)           {}
func (m *mockReporter) RecordSuccess()                  {}
func (m *mockReporter) RecordFailure()                  {}
func (m *mockReporter) RecordSkipped()                  {}

func TestEnricher_EnrichBatch(t *testing.T) {
	repo1 := &domain.ExtractedRepo{RepoID: "owner/repo1", EnrichmentStatus: domain.StatusPending}
	repo2 := &domain.ExtractedRepo{RepoID: "owner/repo2", EnrichmentStatus: domain.StatusPending}
	
	mockRepo := &MockRepo{
		repos: map[string]*domain.ExtractedRepo{
			"owner/repo1": repo1,
			"owner/repo2": repo2,
		},
	}

	mockClient := &MockClient{
		stats: map[string]*domain.RepoStats{
			"owner/repo1": {Stars: 10},
			"owner/repo2": {Stars: 20},
		},
	}

	enricher := NewEnricher(mockRepo, mockClient)

	success, failed, err := enricher.EnrichBatch(context.Background(), 10, false, 2, &mockReporter{})
	if err != nil {
		t.Fatalf("EnrichBatch failed: %v", err)
	}
	if success != 2 {
		t.Errorf("Expected 2 successes, got %d", success)
	}
	if failed != 0 {
		t.Errorf("Expected 0 failures, got %d", failed)
	}
	if *repo1.Stars != 10 {
		t.Errorf("Repo1 stars not updated")
	}
}

func TestEnricher_RateLimit(t *testing.T) {
	repo1 := &domain.ExtractedRepo{RepoID: "owner/limit", EnrichmentStatus: domain.StatusPending}
	mockRepo := &MockRepo{repos: map[string]*domain.ExtractedRepo{"owner/limit": repo1}}
	mockClient := &MockClient{
		failRepo:  "owner/limit",
		failError: domain.ErrRateLimitExceeded,
	}

	enricher := NewEnricher(mockRepo, mockClient)
	_, _, err := enricher.EnrichBatch(context.Background(), 10, false, 1, &mockReporter{})
	if err == nil {
		t.Error("Expected rate limit error, got nil")
	}
	if !errors.Is(err, domain.ErrRateLimitExceeded) {
		t.Errorf("Expected ErrRateLimitExceeded, got %v", err)
	}
}
