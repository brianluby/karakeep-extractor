# Internal Interfaces: Filtering

Update `RankingRepository` to accept a filter string.

```go
type RankingRepository interface {
    // GetRankedRepos returns a list of repos sorted by the criteria and optionally filtered by a tag.
    GetRankedRepos(ctx context.Context, limit int, sortBy RankSortOption, filterTag string) ([]domain.ExtractedRepo, error)
}
```
