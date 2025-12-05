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
	update1 := domain.RepoEnrichmentUpdate{
		RepoID: "owner/python-tool",
		Stats: &domain.RepoStats{
			Stars: 100, Description: "A tool for Python",
		},
		EnrichmentStatus: domain.StatusSuccess,
	}
	repo.Save(ctx, domain.ExtractedRepo{RepoID: update1.RepoID, URL: "url1", Title: "Python Tool", FoundAt: time.Now()})
	repo.UpdateRepoEnrichment(ctx, update1)

	update2 := domain.RepoEnrichmentUpdate{
		RepoID: "owner/go-cli",
		Stats: &domain.RepoStats{
			Stars: 200, Description: "A CLI in Go",
		},
		EnrichmentStatus: domain.StatusSuccess,
	}
	repo.Save(ctx, domain.ExtractedRepo{RepoID: update2.RepoID, URL: "url2", Title: "Go CLI", FoundAt: time.Now()})
	repo.UpdateRepoEnrichment(ctx, update2)

	// Test Filter by Title ("Python")
	repos, err := repo.GetRankedRepos(ctx, 10, domain.SortByStars, "Python")
	if err != nil {
		t.Fatalf("Filter(Python) failed: %v", err)
	}
	if len(repos) != 1 || repos[0].RepoID != "owner/python-tool" {
		t.Errorf("Expected 1 python repo, got %d", len(repos))
	}

	// Test Filter by Description ("CLI")
	repos, err = repo.GetRankedRepos(ctx, 10, domain.SortByStars, "CLI")
	if err != nil {
		t.Fatalf("Filter(CLI) failed: %v", err)
	}
	if len(repos) != 1 || repos[0].RepoID != "owner/go-cli" {
		t.Errorf("Expected 1 go repo, got %d", len(repos))
	}

	// Test Case Insensitivity ("go")
	repos, err = repo.GetRankedRepos(ctx, 10, domain.SortByStars, "go")
	if err != nil {
		t.Fatalf("Filter(go) failed: %v", err)
	}
	if len(repos) != 1 {
		t.Errorf("Expected 1 go repo (case insensitive), got %d", len(repos))
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
