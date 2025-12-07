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

// InitSchema initializes the database schema for ExtractedRepo and Tags.
func (r *SQLiteRepository) InitSchema(ctx context.Context) error {
	const createTableSQL = `
	CREATE TABLE IF NOT EXISTS extracted_repos (
		repo_id TEXT PRIMARY KEY,
		url TEXT NOT NULL,
		source_id TEXT,
		title TEXT,
		found_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := r.db.ExecContext(ctx, createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to initialize schema (extracted_repos): %w", err)
	}

	const createTagsTableSQL = `
	CREATE TABLE IF NOT EXISTS tags (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL
	);
	`
	_, err = r.db.ExecContext(ctx, createTagsTableSQL)
	if err != nil {
		return fmt.Errorf("failed to initialize schema (tags): %w", err)
	}

	const createRepoTagsTableSQL = `
	CREATE TABLE IF NOT EXISTS repo_tags (
		repo_id TEXT NOT NULL,
		tag_id INTEGER NOT NULL,
		PRIMARY KEY (repo_id, tag_id),
		FOREIGN KEY (repo_id) REFERENCES extracted_repos(repo_id) ON DELETE CASCADE,
		FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
	);
	`
	_, err = r.db.ExecContext(ctx, createRepoTagsTableSQL)
	if err != nil {
		return fmt.Errorf("failed to initialize schema (repo_tags): %w", err)
	}

	// Migrations: Add new columns if they don't exist
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
			// fmt.Printf("Migration warning (safe to ignore if column exists): %v\n", err)
		}
	}

	return nil
}

// Save saves an ExtractedRepo to the database. If a repo with the same RepoID already exists, it is ignored.
func (r *SQLiteRepository) Save(ctx context.Context, repo domain.ExtractedRepo) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	const insertRepoSQL = `
	INSERT OR IGNORE INTO extracted_repos (repo_id, url, source_id, title, found_at)
	VALUES (?, ?, ?, ?, ?);
	`
	
	_, err = tx.ExecContext(ctx, insertRepoSQL,
		repo.RepoID,
		repo.URL,
		repo.SourceID,
		repo.Title,
		repo.FoundAt.Format(time.RFC3339),
	)
	if err != nil {
		return fmt.Errorf("failed to save repository: %w", err)
	}

	// Handle Tags
	if len(repo.Tags) > 0 {
		// Note: If repo already exists, we are NOT updating tags here based on "INSERT OR IGNORE" logic above.
		// The requirement says "When extracted again, Then the local database updates to reflect the new tag set."
		// So we probably should update tags even if repo exists.
		// Let's clear existing tags for this repo and re-insert to ensure sync.
		// But wait, if we are using "INSERT OR IGNORE" for repo, we assume it might exist.
		// If it exists, we should ensure it's there.
		// Let's assume 'Save' is called for every extraction.
		
		// Ensure Repo exists (in case it was ignored above but we want to attach tags now)
		// Actually, if it was ignored, the row exists.
		
		// Delete existing links for this repo (full sync strategy)
		const deleteTagsSQL = `DELETE FROM repo_tags WHERE repo_id = ?;
		`
		_, err = tx.ExecContext(ctx, deleteTagsSQL, repo.RepoID)
		if err != nil {
			return fmt.Errorf("failed to clear old tags: %w", err)
		}

		if err := r.saveTags(ctx, tx, repo.RepoID, repo.Tags); err != nil {
			return fmt.Errorf("failed to save tags: %w", err)
		}
	}

	return tx.Commit()
}

// saveTags helper to insert tags and links within a transaction.
func (r *SQLiteRepository) saveTags(ctx context.Context, tx *sql.Tx, repoID string, tags []string) error {
	for _, tag := range tags {
		// 1. Insert Tag (Ignore if exists)
		const insertTagSQL = `INSERT OR IGNORE INTO tags (name) VALUES (?);
		`
		_, err := tx.ExecContext(ctx, insertTagSQL, tag)
		if err != nil {
			return fmt.Errorf("failed to insert tag %s: %w", tag, err)
		}

		// 2. Get Tag ID
		const getTagIDSQL = `SELECT id FROM tags WHERE name = ?;
		`
		var tagID int64
		err = tx.QueryRowContext(ctx, getTagIDSQL, tag).Scan(&tagID)
		if err != nil {
			return fmt.Errorf("failed to get id for tag %s: %w", tag, err)
		}

		// 3. Link Repo to Tag
		const linkTagSQL = `INSERT INTO repo_tags (repo_id, tag_id) VALUES (?, ?);
		`
		_, err = tx.ExecContext(ctx, linkTagSQL, repoID, tagID)
		if err != nil {
			return fmt.Errorf("failed to link tag %s to repo %s: %w", tag, repoID, err)
		}
	}
	return nil
}

// Exists checks if an ExtractedRepo with the given RepoID already exists in the database.
func (r *SQLiteRepository) Exists(ctx context.Context, repoID string) (bool, error) {
	const querySQL = `SELECT COUNT(*) FROM extracted_repos WHERE repo_id = ?;
	`
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
		WHERE repo_id = ?;
		`
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
		WHERE repo_id = ?;
		`
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
		querySQL = `SELECT repo_id, url, source_id, title, found_at, stars, forks, last_pushed_at, description, language, enrichment_status FROM extracted_repos LIMIT ?;
	`
	} else {
		querySQL = `SELECT repo_id, url, source_id, title, found_at, stars, forks, last_pushed_at, description, language, enrichment_status FROM extracted_repos WHERE enrichment_status != 'SUCCESS' LIMIT ?;
	`
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
	baseQuery := `
		SELECT er.repo_id, er.url, er.source_id, er.title, er.found_at, er.stars, er.forks, er.last_pushed_at, er.description, er.language, er.enrichment_status 
		FROM extracted_repos er
		WHERE er.enrichment_status = 'SUCCESS'`

	var args []interface{}
	if filterTag != "" {
		// Old behavior: search title/desc
		// New behavior (from spec): filter by TAGs.
		// "The karakeep rank command MUST support filtering by these locally stored tags via the --tag flag"
		// However, users might expect existing title/desc search to work if --tag is ambiguous.
		// But "filterTag" argument implies tag. Let's implement tag filtering via JOIN.
		
		baseQuery += ` AND EXISTS (
			SELECT 1 FROM repo_tags rt
			JOIN tags t ON rt.tag_id = t.id
			WHERE rt.repo_id = er.repo_id AND t.name = ?
		)`
		args = append(args, filterTag)
	}

	var orderClause string
	switch sortBy {
	case domain.SortByStars:
		orderClause = "ORDER BY er.stars DESC"
	case domain.SortByForks:
		orderClause = "ORDER BY er.forks DESC"
	case domain.SortByUpdated:
		orderClause = "ORDER BY er.last_pushed_at DESC"
	default:
		orderClause = "ORDER BY er.stars DESC"
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