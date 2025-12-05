# Research & Decision Log

## 1. GitHub API Client Library

**Context**: The feature requires fetching repository metadata (stars, forks, etc.) from GitHub.
**Question**: Should we use the official `google/go-github` library or a custom lightweight client using `net/http`?

**Decision**: **Use `net/http` (Standard Library)**

**Rationale**:
*   **Minimizes Dependencies**: The Constitution explicitly states "Minimize external dependencies". A full client library pulls in `go-querystring` and generic structures we don't need.
*   **Simplicity**: We only need a single endpoint (`GET /repos/{owner}/{repo}`) and a few fields. A custom struct and a simple HTTP GET wrapper are sufficient and easier to audit.
*   **Rate Limit Handling**: Parsing `X-RateLimit-*` headers is trivial with standard `http.Response`.

**Alternatives Considered**:
*   `google/go-github`: Robust but heavy. Overkill for just 5 fields.
*   `shurcooL/githubv4` (GraphQL): Reduces over-fetching but adds complexity with GraphQL queries and authentication setup. REST is simpler for this scope.

## 2. SQLite Schema Migration

**Context**: We need to add columns (`stars`, `forks`, `desc`, `pushed_at`, `language`, `enrichment_status`) to the existing `extracted_repos` table.
**Question**: How to apply these schema changes?

**Decision**: **Inline `ALTER TABLE` statements on startup**

**Rationale**:
*   **Consistency**: The project currently initializes tables with `CREATE TABLE IF NOT EXISTS` in `repository.go`.
*   **Simplicity**: Adding columns is safe in SQLite. We can execute `ALTER TABLE extracted_repos ADD COLUMN stars INTEGER` inside a `try/catch` (or check if column exists first) logic block in the repository initialization.
*   **No Extra Tools**: Avoids introducing `golang-migrate` or other heavy migration frameworks for a simple CLI tool.

**Alternatives Considered**:
*   `golang-migrate`: Standard for large apps, but adds file management overhead for a single table change.
*   Re-creating the table: Destructive to user data. Rejected.

## 3. Concurrency Strategy

**Context**: We need to enrich batches of URLs efficiently (FR-004).
**Question**: How to manage concurrency?

**Decision**: **Worker Pool Pattern**

**Rationale**:
*   **Control**: Allows us to set a fixed number of workers (e.g., 5) to respect rate limits and system resources, rather than spawning a goroutine per URL which could flood the network or hit rate limits instantly.
*   **Error Handling**: Easier to aggregate results and errors from a results channel.

**Alternatives Considered**:
*   `sync.WaitGroup` per URL: Harder to throttle.
*   Sequential processing: Too slow (SC-002 requires <10s for 50 repos).

