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

	// Migrations: Add new columns if they don't exist
	// We attempt to add columns and ignore errors (likely "duplicate column name")
	// This is a simple migration strategy for a local CLI tool.
	migrationSQLs := []string{
		`ALTER TABLE extracted_repos ADD COLUMN stars INTEGER;`,
		`ALTER TABLE extracted_repos ADD COLUMN forks INTEGER;`,
		`ALTER TABLE extracted_repos ADD COLUMN last_pushed_at DATETIME;`,
		`ALTER TABLE extracted_repos ADD COLUMN description TEXT;`,
		`ALTER TABLE extracted_repos ADD COLUMN language TEXT;`,
		`ALTER TABLE extracted_repos ADD COLUMN enrichment_status TEXT DEFAULT 'PENDING';`,
	}

	for _, sql := range migrationSQLs {
		if _, err := r.db.ExecContext(ctx, sql); err != nil {
			// Optional: check if err.Error() contains "duplicate column"
			// log.Printf("Migration warning (safe to ignore if column exists): %v", err)
			fmt.Printf("Migration warning (safe to ignore if column exists): %v\n", err)
		}
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

// UpdateRepoEnrichment updates the stats and status of a repository.
func (r *SQLiteRepository) UpdateRepoEnrichment(ctx context.Context, update domain.RepoEnrichmentUpdate) error {
	var updateSQL string
	var args []interface{}

	if update.Stats != nil {
		updateSQL = `
		UPDATE extracted_repos
		SET stars = ?, forks = ?, last_pushed_at = ?, description = ?, language = ?, enrichment_status = ?
		WHERE repo_id = ?;`
		args = []interface{}{
			update.Stats.Stars,
			update.Stats.Forks,
			update.Stats.LastPushed.Format(time.RFC3339),
			update.Stats.Description,
			update.Stats.Language,
			update.EnrichmentStatus,
			update.RepoID,
		}
	} else {
		updateSQL = `
		UPDATE extracted_repos
		SET enrichment_status = ?
		WHERE repo_id = ?;`
		args = []interface{}{
			update.EnrichmentStatus,
			update.RepoID,
		}
	}

	result, err := r.db.ExecContext(ctx, updateSQL, args...)
	if err != nil {
		return fmt.Errorf("failed to update repository enrichment: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("repository not found: %s", update.RepoID)
	}

	return nil
}

// GetReposForEnrichment returns up to 'limit' repos that need enrichment.
func (r *SQLiteRepository) GetReposForEnrichment(ctx context.Context, limit int, force bool) ([]*domain.ExtractedRepo, error) {
	var querySQL string
	if force {
		querySQL = `SELECT repo_id, url, source_id, title, found_at, stars, forks, last_pushed_at, description, language, enrichment_status FROM extracted_repos LIMIT ?;`
	} else {
		querySQL = `SELECT repo_id, url, source_id, title, found_at, stars, forks, last_pushed_at, description, language, enrichment_status FROM extracted_repos WHERE enrichment_status != 'SUCCESS' LIMIT ?;`
	}

	rows, err := r.db.QueryContext(ctx, querySQL, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query repos for enrichment: %w", err)
	}
	defer rows.Close()

	var repos []*domain.ExtractedRepo
	for rows.Next() {
		var r domain.ExtractedRepo
		var foundAt string
		var lastPushedAt sql.NullString // Use NullString for scanning
		var stars, forks sql.NullInt64
		var description, language, enrichmentStatus sql.NullString

		err := rows.Scan(
			&r.RepoID, &r.URL, &r.SourceID, &r.Title, &foundAt,
			&stars, &forks, &lastPushedAt, &description, &language, &enrichmentStatus,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan repo row: %w", err)
		}

		// Parse FoundAt
		t, err := time.Parse(time.RFC3339, foundAt)
		if err != nil {
			// Try parsing as SQLite default CURRENT_TIMESTAMP format "2006-01-02 15:04:05" if RFC3339 fails
			t, err = time.Parse("2006-01-02 15:04:05", foundAt)
			if err != nil {
				return nil, fmt.Errorf("failed to parse found_at time: %w", err)
			}
		}
		r.FoundAt = t

		// Map nullable fields
		if stars.Valid {
			s := int(stars.Int64)
			r.Stars = &s
		}
		if forks.Valid {
			f := int(forks.Int64)
			r.Forks = &f
		}
		if lastPushedAt.Valid {
			t, err := time.Parse(time.RFC3339, lastPushedAt.String)
			if err == nil {
				r.LastPushedAt = &t
			}
		}
		if description.Valid {
			r.Description = &description.String
		}
		if language.Valid {
			r.Language = &language.String
		}
		if enrichmentStatus.Valid {
			r.EnrichmentStatus = domain.EnrichmentStatus(enrichmentStatus.String)
		} else {
			r.EnrichmentStatus = domain.StatusPending // Default if null, though migration sets default
		}

		repos = append(repos, &r)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return repos, nil
}

// GetRankedRepos returns a list of repos sorted by the criteria and optionally filtered by a tag.
func (r *SQLiteRepository) GetRankedRepos(ctx context.Context, limit int, sortBy domain.RankSortOption, filterTag string) ([]domain.ExtractedRepo, error) {
	baseQuery := `SELECT repo_id, url, source_id, title, found_at, stars, forks, last_pushed_at, description, language, enrichment_status 
	              FROM extracted_repos 
	              WHERE enrichment_status = 'SUCCESS'`

	var args []interface{}
	if filterTag != "" {
		baseQuery += " AND (title LIKE ? OR description LIKE ?)"
		likePattern := "%" + filterTag + "%"
		args = append(args, likePattern, likePattern)
	}

	var orderClause string
	switch sortBy {
	case domain.SortByStars:
		orderClause = "ORDER BY stars DESC"
	case domain.SortByForks:
		orderClause = "ORDER BY forks DESC"
	case domain.SortByUpdated:
		orderClause = "ORDER BY last_pushed_at DESC"
	default:
		orderClause = "ORDER BY stars DESC"
	}

	finalQuery := fmt.Sprintf("%s %s LIMIT ?", baseQuery, orderClause)
	args = append(args, limit)

	rows, err := r.db.QueryContext(ctx, finalQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query ranked repos: %w", err)
	}
	defer rows.Close()

	var repos []domain.ExtractedRepo
	for rows.Next() {
		var r domain.ExtractedRepo
		var foundAt string
		var lastPushedAt sql.NullString
		var stars, forks sql.NullInt64
		var description, language, enrichmentStatus sql.NullString

		err := rows.Scan(
			&r.RepoID, &r.URL, &r.SourceID, &r.Title, &foundAt,
			&stars, &forks, &lastPushedAt, &description, &language, &enrichmentStatus,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan repo row: %w", err)
		}

		// Helper to parse time or fallback
		parseTime := func(ts string) time.Time {
			t, err := time.Parse(time.RFC3339, ts)
			if err != nil {
				t, _ = time.Parse("2006-01-02 15:04:05", ts)
			}
			return t
		}

		r.FoundAt = parseTime(foundAt)

		if stars.Valid {
			s := int(stars.Int64)
			r.Stars = &s
		}
		if forks.Valid {
			f := int(forks.Int64)
			r.Forks = &f
		}
		if lastPushedAt.Valid {
			t := parseTime(lastPushedAt.String)
			r.LastPushedAt = &t
		}
		if description.Valid {
			r.Description = &description.String
		}
		if language.Valid {
			r.Language = &language.String
		}
		if enrichmentStatus.Valid {
			r.EnrichmentStatus = domain.EnrichmentStatus(enrichmentStatus.String)
		}

		repos = append(repos, r)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return repos, nil
}
