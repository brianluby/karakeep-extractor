package sqlite

import (
	"context"
	"os"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/brianluby/karakeep-extractor/internal/core/domain"
)

func TestSQLiteRepository_GetRankedRepos_Filter(t *testing.T) {
	db, dbPath := newTestDB(t)
	defer os.Remove(dbPath)
	defer db.Close()

	repo := NewSQLiteRepository(db)
	ctx := context.Background()

	if err := repo.InitSchema(ctx); err != nil {
		t.Fatalf("InitSchema failed: %v", err)
	}

	// Seed data
	repo1 := domain.ExtractedRepo{
		RepoID: "owner/python-tool", 
		URL: "url1", 
		Title: "Python Tool", 
		FoundAt: time.Now(),
		Tags: []string{"python", "data"},
	}
	repo.Save(ctx, repo1)
	repo.UpdateRepoEnrichment(ctx, domain.RepoEnrichmentUpdate{
		RepoID: repo1.RepoID,
		Stats: &domain.RepoStats{Stars: 100},
		EnrichmentStatus: domain.StatusSuccess,
	})

	repo2 := domain.ExtractedRepo{
		RepoID: "owner/go-cli", 
		URL: "url2", 
		Title: "Go CLI", 
		FoundAt: time.Now(),
		Tags: []string{"golang", "cli"},
	}
	repo.Save(ctx, repo2)
	repo.UpdateRepoEnrichment(ctx, domain.RepoEnrichmentUpdate{
		RepoID: repo2.RepoID,
		Stats: &domain.RepoStats{Stars: 200},
		EnrichmentStatus: domain.StatusSuccess,
	})

	// Test Filter by Tag ("python")
	repos, err := repo.GetRankedRepos(ctx, 10, domain.SortByStars, "python")
	if err != nil {
		t.Fatalf("Filter(python) failed: %v", err)
	}
	if len(repos) != 1 || repos[0].RepoID != "owner/python-tool" {
		t.Errorf("Expected 1 python repo, got %d", len(repos))
	}

	// Test Filter by Tag ("cli")
	repos, err = repo.GetRankedRepos(ctx, 10, domain.SortByStars, "cli")
	if err != nil {
		t.Fatalf("Filter(cli) failed: %v", err)
	}
	if len(repos) != 1 || repos[0].RepoID != "owner/go-cli" {
		t.Errorf("Expected 1 go repo, got %d", len(repos))
	}

	// Test No Match
	repos, err = repo.GetRankedRepos(ctx, 10, domain.SortByStars, "java")
	if err != nil {
		t.Fatalf("Filter(java) failed: %v", err)
	}
	if len(repos) != 0 {
		t.Errorf("Expected 0 repos, got %d", len(repos))
	}
}
