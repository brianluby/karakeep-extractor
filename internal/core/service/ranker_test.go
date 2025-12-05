package service_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/brianluby/karakeep-extractor/internal/core/domain"
	"github.com/brianluby/karakeep-extractor/internal/core/service"
)

type mockRankingRepo struct {
	repos []domain.ExtractedRepo
}

func (m *mockRankingRepo) Save(ctx context.Context, repo domain.ExtractedRepo) error { return nil }
func (m *mockRankingRepo) Exists(ctx context.Context, repoID string) (bool, error)   { return false, nil }
func (m *mockRankingRepo) GetReposForEnrichment(ctx context.Context, limit int, force bool) ([]*domain.ExtractedRepo, error) {
	return nil, nil
}
func (m *mockRankingRepo) UpdateRepoEnrichment(ctx context.Context, update domain.RepoEnrichmentUpdate) error {
	return nil
}
func (m *mockRankingRepo) GetRankedRepos(ctx context.Context, limit int, sortBy domain.RankSortOption) ([]domain.ExtractedRepo, error) {
	return m.repos, nil
}

func TestRanker_Rank(t *testing.T) {
	mockRepo := &mockRankingRepo{
		repos: []domain.ExtractedRepo{
			{RepoID: "test/repo1"},
		},
	}
	ranker := service.NewRanker(mockRepo)
	var buf bytes.Buffer

	if err := ranker.Rank(context.Background(), 10, "stars", &buf); err != nil {
		t.Fatalf("Rank failed: %v", err)
	}

	if !strings.Contains(buf.String(), "test/repo1") {
		t.Errorf("Output missing repo: %s", buf.String())
	}
}

func TestRanker_InvalidSort(t *testing.T) {
	ranker := service.NewRanker(&mockRankingRepo{})
	var buf bytes.Buffer
	if err := ranker.Rank(context.Background(), 10, "invalid", &buf); err == nil {
		t.Error("Expected error for invalid sort option")
	}
}
