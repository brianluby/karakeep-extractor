# Data Model: GitHub Repository Ranking

## Entities

### RankedRepoView (Projection)

A simplified view of `ExtractedRepo` specifically for the ranking display.

| Field | Type | Description | Source Field |
|-------|------|-------------|--------------|
| `Rank` | Integer | Position in the list (1-based) | Derived (Row Number) |
| `Name` | String | "owner/repo" | `RepoID` |
| `Stars` | Integer | Star count | `Stars` |
| `Forks` | Integer | Fork count | `Forks` |
| `Updated` | String | Relative time (e.g. "2d ago") | `LastPushedAt` |

## Queries

### Ranking Query

```sql
SELECT 
    repo_id, 
    stars, 
    forks, 
    last_pushed_at 
FROM extracted_repos 
WHERE enrichment_status = 'SUCCESS'
ORDER BY 
    CASE WHEN ? = 'stars' THEN stars END DESC,
    CASE WHEN ? = 'forks' THEN forks END DESC,
    CASE WHEN ? = 'updated' THEN last_pushed_at END DESC
LIMIT ?;
```
*Note: Dynamic ORDER BY in SQLite can be tricky or slow. Constructing the query string safely (allow-list columns) in Go is often better for performance than CASE statements in ORDER BY.*

**Revised Query Strategy (Go string construction):**

```go
baseQuery := "SELECT ... FROM extracted_repos WHERE enrichment_status = 'SUCCESS'"
orderClause := "ORDER BY stars DESC" // Default
if sortBy == "forks" { orderClause = "ORDER BY forks DESC" }
if sortBy == "updated" { orderClause = "ORDER BY last_pushed_at DESC" }

finalQuery := fmt.Sprintf("%s %s LIMIT ?", baseQuery, orderClause)
```
