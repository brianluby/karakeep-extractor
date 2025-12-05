# CLI Contract

## Command: `extract`

**Description**: Fetches bookmarks from Karakeep, filters for GitHub repositories, and saves them to SQLite.

### Usage
```bash
karakeep-extractor extract [flags]
```

### Flags
| Flag | Env Var | Description | Default |
|------|---------|-------------|---------|
| `--url` | `KARAKEEP_URL` | Base URL of the Karakeep instance | (Required) |
| `--token` | `KARAKEEP_TOKEN` | API Bearer Token | (Required) |
| `--db` | `KARAKEEP_DB` | Path to SQLite database | `./karakeep.db` |

### Exit Codes
*   `0`: Success (Extraction completed).
*   `1`: Configuration Error (Missing flags/env vars).
*   `2`: Connection Error (Network/Auth failure).
*   `3`: Runtime Error (DB write failure, etc.).

### Output (Stdout)
*   Progress logs (e.g., "Fetched page 1...", "Found 5 repos...").
*   Final summary ("Total repos extracted: X").

### Output (Stderr)
*   Error logs (e.g., "Failed to connect...", "Skipping malformed URL...").
