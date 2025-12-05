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
