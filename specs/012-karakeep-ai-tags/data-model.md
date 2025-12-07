# Data Model: Karakeep AI Tags Extraction

## Domain Entities

### RawBookmark (Updated)
Updated to include the tags field from the API.

```go
type RawBookmark struct {
	ID      string  `json:"id"`
	Title   *string `json:"title"`
	Content struct {
		URL         string `json:"url"`
		Title       string `json:"title"`
		Description string `json:"description"`
		HTMLContent string `json:"htmlContent"`
	} `json:"content"`
    // New Field
    Tags []string `json:"tags"` 
}
```

### ExtractedRepo (Updated)
Updated to carry tags through the system.

```go
type ExtractedRepo struct {
	// ... existing fields ...
    
    // New Field
    Tags []string
}
```

## Database Schema

### New Table: `tags`
| Column | Type | Constraints |
|---|---|---|
| `id` | INTEGER | PRIMARY KEY AUTOINCREMENT |
| `name` | TEXT | UNIQUE, NOT NULL |

### New Table: `repo_tags`
| Column | Type | Constraints |
|---|---|---|
| `repo_id` | TEXT | NOT NULL, FOREIGN KEY -> `extracted_repos.repo_id` |
| `tag_id` | INTEGER | NOT NULL, FOREIGN KEY -> `tags.id` |
| **PK** | | `PRIMARY KEY (repo_id, tag_id)` |

## Migration Strategy
The application will run a migration on startup to create these tables if they don't exist.
