# Data Model: Filtering

## Queries

### Filtered Ranking Query

```sql
SELECT 
    repo_id, 
    stars, 
    forks, 
    last_pushed_at,
    title,
    description
FROM extracted_repos 
WHERE enrichment_status = 'SUCCESS'
  AND (
      title LIKE ? 
      OR description LIKE ?
  )
ORDER BY ...
LIMIT ?;
```

**Parameters**:
- `?`: `"%tag%"` (case-insensitive by default in SQLite for ASCII).

## Entities

No changes to Go structs (`ExtractedRepo`) required, just the query logic.
