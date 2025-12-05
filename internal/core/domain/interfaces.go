package domain

import "context"

// BookmarkSource Port: Source (Secondary)
type BookmarkSource interface {
	FetchBookmarks(ctx context.Context, page int) ([]RawBookmark, error)
}

// RepoEnrichmentUpdate defines the fields to update for an enrichment operation.
type RepoEnrichmentUpdate struct {
	RepoID           string
	Stats            *RepoStats // Can be nil if only status updates
	EnrichmentStatus EnrichmentStatus
}

// RepoRepository Port: Storage (Secondary)
type RepoRepository interface {
	Save(ctx context.Context, repo ExtractedRepo) error
	Exists(ctx context.Context, repoID string) (bool, error)

	// GetReposForEnrichment returns up to 'limit' repos that need enrichment.
	// If force is true, returns any repo. If false, only those with EnrichmentStatus != SUCCESS.
	GetReposForEnrichment(ctx context.Context, limit int, force bool) ([]*ExtractedRepo, error)

	// UpdateRepoEnrichment updates the stats and status of a repository.
	UpdateRepoEnrichment(ctx context.Context, update RepoEnrichmentUpdate) error
}

// GitHubClient Port: Source (Secondary)
type GitHubClient interface {
	// GetRepoStats fetches metadata for a single repo.
	// Returns stats, remaining rate limit, and error.
	// Returns specific error for 404.
	GetRepoStats(ctx context.Context, owner, repo string) (*RepoStats, int, error)
}
