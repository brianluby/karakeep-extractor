# Research: Karakeep Data Extraction

**Feature**: Karakeep Data Extraction
**Status**: Phase 0 Complete

## Technical Decisions

### 1. HTTP Client & Retry Logic
*   **Decision**: Use Go standard `net/http` with a custom `RoundTripper` or middleware for exponential backoff.
*   **Rationale**: Avoids heavy 3rd party libraries for simple HTTP calls. Standard library is robust.
*   **Alternatives**: `go-resty` (rejected: extra dependency).

### 2. SQLite Driver
*   **Decision**: `modernc.org/sqlite` (Pure Go) or `github.com/mattn/go-sqlite3` (CGO).
*   **Selected**: `github.com/mattn/go-sqlite3`.
*   **Rationale**: Industry standard, highly reliable.
*   **Tradeoff**: Requires CGO, which might complicate cross-compilation, but acceptable for this stage.

### 3. CLI Argument Parsing
*   **Decision**: Go standard `flag` package with `os.Args` handling for subcommands (e.g., `extract`).
*   **Rationale**: Keeps the binary small and follows "Minimize external dependencies" constitution principle.
*   **Alternatives**: `cobra` (rejected: overkill for single-purpose tool currently).

### 4. URL Normalization
*   **Decision**: Use `net/url` parser.
*   **Logic**:
    1. Parse URL.
    2. Check Host == `github.com`.
    3. Trim path to first 2 segments (`owner/repo`).
    4. Lowercase.
    5. Discard fragments/queries.

## Unknowns & Resolutions

*   **Karakeep API Endpoint**: Assumed to be RESTful.
    *   *Resolution*: Configurable via `KARAKEEP_URL`. The code will append path if needed or expect full URL. Defaulting to expecting base URL and appending `/api/bookmarks`.

## Constitution Compliance
*   **CLI First**: Yes.
*   **Modular**: `internal/adapter/karakeep` will handle API, `internal/core/service` will handle logic.
*   **Dependencies**: Kept to minimum (only SQLite driver).
