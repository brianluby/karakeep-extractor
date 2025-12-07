# Data Model: Advanced Local Database Querying

## Domain Entities

### QueryFilter
Represents the criteria for filtering repositories.

```go
type QueryFilter struct {
    MinStars     *int
    MaxStars     *int
    MinForks     *int
    MaxForks     *int
    CreatedAfter *time.Time // Maps to FoundAt > X
    CreatedBefore *time.Time // Maps to FoundAt < X
    Language     *string
    Tag          *string
    SearchTerm   *string
}
```

## Database Interface

### New Method: `QueryRepos`

```go
QueryRepos(ctx context.Context, filter QueryFilter) ([]ExtractedRepo, error)
```

This method will replace or augment the existing `GetRankedRepos` if the filters become too complex for the old signature. Ideally, `GetRankedRepos` logic can be refactored to use `QueryFilter`.
