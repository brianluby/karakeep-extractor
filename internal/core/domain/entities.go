package domain

import "time"

// KarakeepConfig Configuration for connecting to the source.
type KarakeepConfig struct {
	BaseURL string
	APIToken string
}

// RawBookmark Represents a bookmark as returned by the Karakeep API.
type RawBookmark struct {
	ID      string `json:"id"`
	Title   *string `json:"title"` // Top-level title can be null
	Content struct {
		URL         string `json:"url"`
		Title       string `json:"title"`
		Description string `json:"description"`
		HTMLContent string `json:"htmlContent"`
	} `json:"content"`
}

type EnrichmentStatus string

const (
	StatusPending  EnrichmentStatus = "PENDING"
	StatusSuccess  EnrichmentStatus = "SUCCESS"
	StatusNotFound EnrichmentStatus = "NOT_FOUND"
	StatusAPIError EnrichmentStatus = "API_ERROR"
)

// RepoStats represents the metadata fetched from GitHub
type RepoStats struct {
	Stars       int
	Forks       int
	LastPushed  time.Time
	Description string
	Language    string
}

// ExtractedRepo The refined domain entity representing a GitHub repository found in bookmarks.
type ExtractedRepo struct {
	RepoID   string    // Canonical "owner/name" (Primary Key in DB).
	URL      string    // Normalized HTTPS URL.
	SourceID string    // ID of the original Karakeep bookmark.
	Title    string    // Title from the bookmark.
	FoundAt  time.Time // Timestamp of extraction.

	// Enrichment Data
	Stars            *int             // Nullable
	Forks            *int             // Nullable
	LastPushedAt     *time.Time       // Nullable
	Description      *string          // Nullable
	Language         *string          // Nullable
	EnrichmentStatus EnrichmentStatus
}