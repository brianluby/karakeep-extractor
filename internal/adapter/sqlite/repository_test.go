package sqlite

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/brianluby/karakeep-extractor/internal/core/domain"
)

func newTestDB(t *testing.T) (*sql.DB, string) {
	f, err := os.CreateTemp("", "testdb_*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	dbPath := f.Name()
	f.Close()

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	return db, dbPath
}

func TestSQLiteRepository_InitSchema(t *testing.T) {
	db, dbPath := newTestDB(t)
	defer os.Remove(dbPath)
	defer db.Close()

	repo := NewSQLiteRepository(db)
	err := repo.InitSchema(context.Background())
	if err != nil {
		t.Fatalf("InitSchema failed: %v", err)
	}

	// Verify table exists
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name='extracted_repos';")
	if err != nil {
		t.Fatalf("Failed to query sqlite_master: %v", err)
	}
	defer rows.Close()

	if !rows.Next() {
		t.Error("Table 'extracted_repos' not created")
	}
}

func TestSQLiteRepository_SaveAndExists(t *testing.T) {
	db, dbPath := newTestDB(t)
	defer os.Remove(dbPath)
	defer db.Close()

	repo := NewSQLiteRepository(db)
	err := repo.InitSchema(context.Background())
	if err != nil {
		t.Fatalf("InitSchema failed: %v", err)
	}

	ctx := context.Background()
	testRepo := domain.ExtractedRepo{
		RepoID:   "owner/repo",
		URL:      "https://github.com/owner/repo",
		SourceID: "src123",
		Title:    "Test Repo",
		FoundAt:  time.Now().Truncate(time.Second), // Truncate to match SQLite precision
	}

	// Test Save
	err = repo.Save(ctx, testRepo)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Verify Exists
	exists, err := repo.Exists(ctx, "owner/repo")
	if err != nil {
		t.Fatalf("Exists failed: %v", err)
	}
	if !exists {
		t.Error("Expected repo to exist, but it doesn't")
	}

	exists, err = repo.Exists(ctx, "nonexistent/repo")
	if err != nil {
		t.Fatalf("Exists failed: %v", err)
	}
	if exists {
		t.Error("Expected repo not to exist, but it does")
	}

	// Test saving a duplicate (should not error, but exists should be true)
	err = repo.Save(ctx, testRepo)
	if err != nil {
		t.Fatalf("Save duplicate failed: %v", err) // Our repo should handle it gracefully, not error
	}

	// Verify saved data
	var (
		repoID   string
		url      string
		sourceID string
		title    string
		foundAt  time.Time
	)
	row := db.QueryRow("SELECT repo_id, url, source_id, title, found_at FROM extracted_repos WHERE repo_id = ?", "owner/repo")
	err = row.Scan(&repoID, &url, &sourceID, &title, &foundAt)
	if err != nil {
		t.Fatalf("Failed to scan saved repo: %v", err)
	}

	if repoID != testRepo.RepoID || url != testRepo.URL || sourceID != testRepo.SourceID || title != testRepo.Title || !foundAt.Equal(testRepo.FoundAt) {
		t.Errorf("Retrieved repo mismatch. Expected %+v, Got repo_id:%s, url:%s, source_id:%s, title:%s, found_at:%s", testRepo, repoID, url, sourceID, title, foundAt.String())
	}
}

func TestSQLiteRepository_Enrichment(t *testing.T) {
	db, dbPath := newTestDB(t)
	defer os.Remove(dbPath)
	defer db.Close()

	repo := NewSQLiteRepository(db)
	ctx := context.Background()

	if err := repo.InitSchema(ctx); err != nil {
		t.Fatalf("InitSchema failed: %v", err)
	}

	// Seed data
	seeds := []domain.ExtractedRepo{
		{RepoID: "owner/repo1", URL: "https://github.com/owner/repo1", FoundAt: time.Now()},
		{RepoID: "owner/repo2", URL: "https://github.com/owner/repo2", FoundAt: time.Now()},
	}
	for _, s := range seeds {
		if err := repo.Save(ctx, s); err != nil {
			t.Fatalf("Failed to seed repo %s: %v", s.RepoID, err)
		}
	}

	// Test GetReposForEnrichment (Should return all initially as pending)
	repos, err := repo.GetReposForEnrichment(ctx, 10, false)
	if err != nil {
		t.Fatalf("GetReposForEnrichment failed: %v", err)
	}
	if len(repos) != 2 {
		t.Errorf("Expected 2 repos for enrichment, got %d", len(repos))
	}

	// Test UpdateRepoEnrichment
	stars := 100
	desc := "Test Description"
	update := domain.RepoEnrichmentUpdate{
		RepoID: "owner/repo1",
		Stats: &domain.RepoStats{
			Stars:       stars,
			Description: desc,
			LastPushed:  time.Now(),
			Language:    "Go",
			Forks:       10,
		},
		EnrichmentStatus: domain.StatusSuccess,
	}

	if err := repo.UpdateRepoEnrichment(ctx, update); err != nil {
		t.Fatalf("UpdateRepoEnrichment failed: %v", err)
	}

	// Test GetReposForEnrichment again (Should return only repo2)
	repos, err = repo.GetReposForEnrichment(ctx, 10, false)
	if err != nil {
		t.Fatalf("GetReposForEnrichment failed: %v", err)
	}
	if len(repos) != 1 {
		t.Errorf("Expected 1 repo for enrichment, got %d", len(repos))
	}
	if repos[0].RepoID != "owner/repo2" {
		t.Errorf("Expected repo2, got %s", repos[0].RepoID)
	}

	// Verify Update persisted
	// Use GetReposForEnrichment with force=true to fetch enriched repo
	// Or query direct. Let's rely on GetReposForEnrichment force logic if possible or just check force flag.
	repos, err = repo.GetReposForEnrichment(ctx, 10, true) // Force get all
	if err != nil {
		t.Fatalf("GetReposForEnrichment force failed: %v", err)
	}
	
	var enrichedRepo *domain.ExtractedRepo
	for _, r := range repos {
		if r.RepoID == "owner/repo1" {
			enrichedRepo = r
			break
		}
	}

	if enrichedRepo == nil {
		t.Fatal("Enriched repo not found in results")
	}
	if enrichedRepo.Stars == nil || *enrichedRepo.Stars != stars {
		t.Errorf("Expected stars %d, got %v", stars, enrichedRepo.Stars)
	}
	if enrichedRepo.Description == nil || *enrichedRepo.Description != desc {
		t.Errorf("Expected description %s, got %v", desc, enrichedRepo.Description)
	}
	if enrichedRepo.EnrichmentStatus != domain.StatusSuccess {
		t.Errorf("Expected status SUCCESS, got %s", enrichedRepo.EnrichmentStatus)
	}
}
