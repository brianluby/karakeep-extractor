package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
	"github.com/brianluby/karakeep-extractor/internal/core/domain"
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{db: db}
}

// InitSchema initializes the database schema for ExtractedRepo.
func (r *SQLiteRepository) InitSchema(ctx context.Context) error {
	const createTableSQL = `
	CREATE TABLE IF NOT EXISTS extracted_repos (
		repo_id TEXT PRIMARY KEY,
		url TEXT NOT NULL,
		source_id TEXT,
		title TEXT,
		found_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := r.db.ExecContext(ctx, createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to initialize schema: %w", err)
	}
	return nil
}

// Save saves an ExtractedRepo to the database. If a repo with the same RepoID already exists, it is ignored.
func (r *SQLiteRepository) Save(ctx context.Context, repo domain.ExtractedRepo) error {
	const insertSQL = `
	INSERT OR IGNORE INTO extracted_repos (repo_id, url, source_id, title, found_at)
	VALUES (?, ?, ?, ?, ?);`
	
	_, err := r.db.ExecContext(ctx, insertSQL,
		repo.RepoID,
		repo.URL,
		repo.SourceID,
		repo.Title,
		repo.FoundAt.Format(time.RFC3339), // Store as ISO 8601 string
	)
	if err != nil {
		return fmt.Errorf("failed to save repository: %w", err)
	}
	return nil
}

// Exists checks if an ExtractedRepo with the given RepoID already exists in the database.
func (r *SQLiteRepository) Exists(ctx context.Context, repoID string) (bool, error) {
	const querySQL = `SELECT COUNT(*) FROM extracted_repos WHERE repo_id = ?;`
	var count int
	err := r.db.QueryRowContext(ctx, querySQL, repoID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check if repository exists: %w", err)
	}
	return count > 0, nil
}
