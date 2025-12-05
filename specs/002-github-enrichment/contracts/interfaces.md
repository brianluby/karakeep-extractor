# Internal Interfaces

## GitHub Client

Adapter layer for interacting with GitHub API.

```go
package port

import (
    "context"
    "time"
)

// RepoStats represents the data fetched from GitHub
type RepoStats struct {
    Stars       int
    Forks       int
    LastPushed  time.Time
    Description string
    Language    string
}

type GitHubClient interface {
    // GetRepoStats fetches metadata for a single repo.
    // Returns stats, remaining rate limit, and error.
    // Returns specific error for 404.
    GetRepoStats(ctx context.Context, owner, repo string) (*RepoStats, int, error)
}
```

## Repository (Storage)

Extensions to the existing storage interface.

```go
package port

import "context"

type EnrichmentStatus string

const (
    StatusPending  EnrichmentStatus = "PENDING"
    StatusSuccess  EnrichmentStatus = "SUCCESS"
    StatusNotFound EnrichmentStatus = "NOT_FOUND"
    StatusAPIError EnrichmentStatus = "API_ERROR"
)

type RepoEnrichmentUpdate struct {
    ID               int
    Stats            *RepoStats // Can be nil if only status updates
    EnrichmentStatus EnrichmentStatus
}

type Repository interface {
    // GetReposForEnrichment returns up to 'limit' repos that need enrichment.
    // If force is true, returns any repo. If false, only those with EnrichmentStatus != SUCCESS.
    GetReposForEnrichment(ctx context.Context, limit int, force bool) ([]*Entity, error)

    // UpdateRepoEnrichment updates the stats and status of a repository.
    UpdateRepoEnrichment(ctx context.Context, update RepoEnrichmentUpdate) error
}
```
