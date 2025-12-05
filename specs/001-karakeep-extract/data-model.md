# Data Model: Karakeep Extraction

## Domain Entities

### 1. KarakeepConfig
Configuration for connecting to the source.
*   `BaseURL` (string): URL of the Karakeep instance.
*   `APIToken` (string): Bearer token for authentication.

### 2. RawBookmark
Represents a bookmark as returned by the Karakeep API.
*   `ID` (string/int): Unique identifier from Karakeep.
*   `URL` (string): The saved URL.
*   `Title` (string): Title of the bookmark.
*   `Content` (string): Description or summary content (may contain links).

### 3. ExtractedRepo
The refined domain entity representing a GitHub repository found in bookmarks.
*   `RepoID` (string): Canonical "owner/name" (Primary Key in DB).
*   `URL` (string): Normalized HTTPS URL.
*   `SourceID` (string): ID of the original Karakeep bookmark.
*   `Title` (string): Title from the bookmark.
*   `FoundAt` (datetime): Timestamp of extraction.

## Database Schema (SQLite)

```sql
CREATE TABLE IF NOT EXISTS extracted_repos (
    repo_id TEXT PRIMARY KEY, -- "owner/name"
    url TEXT NOT NULL,
    source_id TEXT,
    title TEXT,
    found_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

## Interfaces (Go)

### Port: Source (Secondary)
```go
type BookmarkSource interface {
    FetchBookmarks(ctx context.Context, page int) ([]RawBookmark, error)
}
```

### Port: Storage (Secondary)
```go
type RepoRepository interface {
    Save(ctx context.Context, repo ExtractedRepo) error
    Exists(ctx context.Context, repoID string) (bool, error)
}
```
