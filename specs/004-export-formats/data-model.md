# Data Model: Export Formats

## Entities

### RankedRepo (JSON Representation)

When exporting to JSON, the structure should be explicit.

```json
[
  {
    "rank": 1,
    "repo_id": "owner/repo",
    "url": "https://github.com/owner/repo",
    "stars": 1500,
    "forks": 200,
    "last_pushed_at": "2023-10-27T10:00:00Z",
    "description": "A cool library",
    "language": "Go"
  }
]
```

### RankedRepo (CSV Representation)

```csv
rank,repo_id,url,stars,forks,last_pushed_at,description,language
1,owner/repo,https://github.com/owner/repo,1500,200,2023-10-27T10:00:00Z,"A cool library",Go
```

## Interfaces

### Exporter

```go
type Exporter interface {
    Export(repos []domain.ExtractedRepo, writer io.Writer) error
}
```

### Sink

```go
type Sink interface {
    Send(ctx context.Context, repos []domain.ExtractedRepo) error
}
```
