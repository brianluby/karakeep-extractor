package domain

import "time"

// KarakeepConfig Configuration for connecting to the source.
type KarakeepConfig struct {
	BaseURL string
	APIToken string
}

// RawBookmark Represents a bookmark as returned by the Karakeep API.
type RawBookmark struct {
	ID      string `json:"id"` // Karakeep ID can be int or string, using string for flexibility
	URL     string `json:"url"`
	Title   string `json:"title"`
	Content string `json:"content"` // Description or summary content (may contain links).
}

// ExtractedRepo The refined domain entity representing a GitHub repository found in bookmarks.
type ExtractedRepo struct {
	RepoID   string    // Canonical "owner/name" (Primary Key in DB).
	URL      string    // Normalized HTTPS URL.
	SourceID string    // ID of the original Karakeep bookmark.
	Title    string    // Title from the bookmark.
	FoundAt  time.Time // Timestamp of extraction.
}