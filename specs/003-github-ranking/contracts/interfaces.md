# Internal Interfaces

## Ranking Repository (Port)

Extension to `RepoRepository` or a new specialized interface.

```go
package port

import (
    "context"
    "github.com/brianluby/karakeep-extractor/internal/core/domain"
)

type RankSortOption string

const (
    SortByStars   RankSortOption = "stars"
    SortByForks   RankSortOption = "forks"
    SortByUpdated RankSortOption = "updated"
)

type RankingRepository interface {
    // GetRankedRepos returns a list of repos sorted by the criteria.
    GetRankedRepos(ctx context.Context, limit int, sortBy RankSortOption) ([]domain.ExtractedRepo, error)
}
```
