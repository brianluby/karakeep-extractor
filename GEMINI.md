# karakeep-extractor Development Guidelines

Auto-generated from all feature plans. Last updated: 2025-12-04

## Active Technologies
- Go 1.25.5 + `net/http` (Standard Library) for API interaction; `encoding/json` for parsing; `github.com/mattn/go-sqlite3` (Existing) for storage. (002-github-enrichment)
- SQLite (Extension of existing `extracted_repos` table). (002-github-enrichment)
- Go 1.25.5 + `net/http` (Stdlib), `encoding/json` (Stdlib), `github.com/mattn/go-sqlite3` (Existing) (003-github-ranking)
- SQLite (Read-only access for this feature) (003-github-ranking)
- Go 1.25.5 + `net/http` (Stdlib) for Sink POST; `encoding/json` (Stdlib) for JSON export; `encoding/csv` (Stdlib) for CSV export. (004-export-formats)
- SQLite (Read-only access via existing `RankingRepository`). (004-export-formats)
- Go 1.25.5 + `gopkg.in/yaml.v3` (for YAML parsing - standard in Go ecosystem for config), `bufio` (Stdlib) for interactive prompts. (005-configuration-wizard)
- YAML file (`~/.config/karakeep/config.yaml`). (005-configuration-wizard)
- Go 1.25.5 + `net/http` (Stdlib) for API calls. (006-trillium-integration)
- Configuration update (YAML) to store Trillium credentials. (006-trillium-integration)
- Go 1.25.5 + `github.com/mattn/go-sqlite3` (Existing). (007-filter-by-tag)
- SQLite (Read-only query modification). (007-filter-by-tag)
- Markdown (for docs), Go 1.25.5 (for CLI help text updates). + None. (008-documentation-sprint)

- Go 1.25.5 + `net/http` (Stdlib), `github.com/mattn/go-sqlite3` (SQLite Driver) (001-karakeep-extract)

## Project Structure

```text
src/
tests/
```

## Commands

# Add commands for Go 1.25.5

## Code Style

Go 1.25.5: Follow standard conventions

## Recent Changes
- 008-documentation-sprint: Added Markdown (for docs), Go 1.25.5 (for CLI help text updates). + None.
- 007-filter-by-tag: Added Go 1.25.5 + `github.com/mattn/go-sqlite3` (Existing).
- 006-trillium-integration: Added Go 1.25.5 + `net/http` (Stdlib) for API calls.


<!-- MANUAL ADDITIONS START -->
<!-- MANUAL ADDITIONS END -->
