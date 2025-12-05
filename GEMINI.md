# karakeep-extractor Development Guidelines

Auto-generated from all feature plans. Last updated: 2025-12-04

## Active Technologies
- Go 1.25.5 + `net/http` (Standard Library) for API interaction; `encoding/json` for parsing; `github.com/mattn/go-sqlite3` (Existing) for storage. (002-github-enrichment)
- SQLite (Extension of existing `extracted_repos` table). (002-github-enrichment)
- Go 1.25.5 + `net/http` (Stdlib), `encoding/json` (Stdlib), `github.com/mattn/go-sqlite3` (Existing) (003-github-ranking)
- SQLite (Read-only access for this feature) (003-github-ranking)
- Go 1.25.5 + `net/http` (Stdlib) for Sink POST; `encoding/json` (Stdlib) for JSON export; `encoding/csv` (Stdlib) for CSV export. (004-export-formats)
- SQLite (Read-only access via existing `RankingRepository`). (004-export-formats)

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
- 004-export-formats: Added Go 1.25.5 + `net/http` (Stdlib) for Sink POST; `encoding/json` (Stdlib) for JSON export; `encoding/csv` (Stdlib) for CSV export.
- 003-github-ranking: Added Go 1.25.5 + `net/http` (Stdlib), `encoding/json` (Stdlib), `github.com/mattn/go-sqlite3` (Existing)
- 002-github-enrichment: Added Go 1.25.5 + `net/http` (Standard Library) for API interaction; `encoding/json` for parsing; `github.com/mattn/go-sqlite3` (Existing) for storage.


<!-- MANUAL ADDITIONS START -->
<!-- MANUAL ADDITIONS END -->
