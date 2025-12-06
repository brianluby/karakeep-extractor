package domain

import (
	"context"
)

// BookmarkSource Interface for fetching bookmarks.
type BookmarkSource interface {
	FetchBookmarks(ctx context.Context, page int) ([]RawBookmark, error)
}

// RepoRepository Interface for persisting extracted repositories.
type RepoRepository interface {
	Save(ctx context.Context, repo ExtractedRepo) error
	Exists(ctx context.Context, repoID string) (bool, error)
	GetReposForEnrichment(ctx context.Context, limit int, force bool) ([]*ExtractedRepo, error)
	UpdateRepoEnrichment(ctx context.Context, update RepoEnrichmentUpdate) error
}

type RepoEnrichmentUpdate struct {
	RepoID           string
	Stats            *RepoStats
	EnrichmentStatus EnrichmentStatus
}

// GitHubClient Interface for fetching metadata from GitHub.
type GitHubClient interface {
	GetRepoStats(ctx context.Context, owner, name string) (*RepoStats, int, error)
}

// RankingRepository interface for querying ranked repos (ReadOnly usually)
type RankingRepository interface {
	GetRankedRepos(ctx context.Context, limit int, sortBy RankSortOption, tagFilter string) ([]ExtractedRepo, error)
}

type RankSortOption string

const (
	SortByStars   RankSortOption = "stars"
	SortByForks   RankSortOption = "forks"
	SortByUpdated RankSortOption = "updated"
)

// Sink interface for exporting data to external services
type Sink interface {
	Send(ctx context.Context, repos []ExtractedRepo) error
}

// ProgressReporter abstracts the output mechanism (CLI logs vs TUI updates).
type ProgressReporter interface {
	// Start initializes the progress tracking.
	// total: expected number of items (-1 if unknown).
	// title: description of the task.
	Start(total int, title string)

	// Increment adds to the processed count.
	Increment()

	// SetStatus updates the description of the current item being processed.
	SetStatus(status string)

	// Log records a message (info/error) without stopping the process.
	// In TUI mode, this goes to the log tail. In text mode, this is stderr/stdout.
	Log(message string)

	// Error records an error specifically (may be highlighted differently).
	Error(err error)

	// Finish signals completion with a summary message.
	Finish(summary string)
}
