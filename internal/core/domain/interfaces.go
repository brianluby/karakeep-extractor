package domain

import "context"

// BookmarkSource Port: Source (Secondary)
type BookmarkSource interface {
	FetchBookmarks(ctx context.Context, page int) ([]RawBookmark, error)
}

// RepoRepository Port: Storage (Secondary)
type RepoRepository interface {
	Save(ctx context.Context, repo ExtractedRepo) error
	Exists(ctx context.Context, repoID string) (bool, error)
}
