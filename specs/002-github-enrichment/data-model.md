# Data Model: GitHub Enrichment

## Entities

### RepoStats (Embeddable)

Represents the metadata fetched from GitHub.

| Field | Type | Description |
|-------|------|-------------|
| `stars` | Integer | Number of stargazers (`stargazers_count`) |
| `forks` | Integer | Number of forks (`forks_count`) |
| `last_pushed_at` | DateTime | Timestamp of last push (`pushed_at`) |
| `description` | String | Repository description |
| `language` | String | Primary coding language |

### EnrichedRepo (Update to ExtractedRepo)

The existing `extracted_repos` table will be expanded.

| Field | Type | Constraint | New? | Description |
|-------|------|------------|------|-------------|
| `id` | INTEGER | PK | No | Existing ID |
| `url` | TEXT | UNIQUE | No | GitHub URL |
| `status` | TEXT | | No | Extraction status |
| `created_at` | DATETIME| | No | |
| `stars` | INTEGER | NULLable | **Yes** | Enriched star count |
| `forks` | INTEGER | NULLable | **Yes** | Enriched fork count |
| `last_pushed_at` | DATETIME| NULLable | **Yes** | GitHub pushed_at |
| `description` | TEXT | NULLable | **Yes** | |
| `language` | TEXT | NULLable | **Yes** | |
| `enrichment_status`| TEXT | | **Yes** | Enum: `PENDING`, `SUCCESS`, `NOT_FOUND`, `API_ERROR` |

## Schema Updates (SQLite)

```sql
ALTER TABLE extracted_repos ADD COLUMN stars INTEGER;
ALTER TABLE extracted_repos ADD COLUMN forks INTEGER;
ALTER TABLE extracted_repos ADD COLUMN last_pushed_at DATETIME;
ALTER TABLE extracted_repos ADD COLUMN description TEXT;
ALTER TABLE extracted_repos ADD COLUMN language TEXT;
ALTER TABLE extracted_repos ADD COLUMN enrichment_status TEXT DEFAULT 'PENDING';
```
*Note: In implementation, check if columns exist before adding to avoid errors on multiple runs.*
