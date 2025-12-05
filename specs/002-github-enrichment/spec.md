# Feature Specification: GitHub Repository Enrichment

**Feature Branch**: `002-github-enrichment`
**Created**: 2025-12-04
**Status**: Draft
**Input**: User description: "Implement GitHub API client for repository enrichment"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Connect and Authenticate (Priority: P1)

As a user, I want the tool to connect to the GitHub API, optionally using my personal access token, so that I can fetch repository data without hitting strict rate limits.

**Why this priority**: Essential for reliable data fetching; unauthenticated limits are too low for typical use cases (60/hr vs 5000/hr).

**Independent Test**: Can be tested by providing a valid token and verifying the rate limit quota, vs providing no token and seeing the lower quota.

**Acceptance Scenarios**:

1. **Given** a valid GitHub Personal Access Token (PAT), **When** the tool initializes, **Then** it should authenticate and log the available rate limit (e.g., 5000 requests remaining).
2. **Given** no token is provided, **When** the tool initializes, **Then** it should proceed in unauthenticated mode but warn the user about lower rate limits.
3. **Given** an invalid token, **When** the tool initializes, **Then** it should report an authentication error and exit (or fall back to unauthenticated if configured, but explicit failure is safer for clarity).

---

### User Story 2 - Enrich Repository Data (Priority: P1)

As a user, I want the tool to scan my extraction database for GitHub URLs and fetch their details (stars, forks, last update, description), so I can populate my local data with popularity and activity metrics.

**Why this priority**: This is the core value proposition—turning a raw URL into actionable data.

**Independent Test**: Can be tested by seeding the database with known public repository URLs and verifying the rows are updated with live data.

**Acceptance Scenarios**:

1. **Given** the database contains valid GitHub repository URLs, **When** the enrichment process runs, **Then** it should update each row with: Star Count, Fork Count, Last Pushed Date, Description, and Primary Language.
2. **Given** a database row with a URL that has been deleted or made private (404), **When** enriched, **Then** its `enrichment_status` column should be set to "NOT_FOUND" or "INACCESSIBLE".
3. **Given** a URL that is not a repo (e.g., a user profile), **When** enriched, **Then** it should be handled gracefully (skipped or flagged).
---

### Edge Cases

- **Rate Limit Exceeded**: If the batch size exceeds the remaining API quota, the system should process what it can, save current progress, and exit gracefully with an informative message.
- **Moved Repositories**: GitHub redirects moved repos. The client should follow redirects or report the new location.
- **Network Flakiness**: Transient network errors should be retried automatically (e.g., exponential backoff).

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST accept a GitHub Personal Access Token (PAT) via configuration or environment variable.
- **FR-002**: System MUST fetch repository metadata from the GitHub REST API, specifically: `stargazers_count`, `forks_count`, `pushed_at`, `description`, `language`.
- **FR-003**: System MUST handle API Rate Limiting by checking headers (`X-RateLimit-Remaining`). If limits are hit, it MUST save current progress and exit gracefully with an informative message.
- **FR-004**: System MUST process a batch of URLs efficiently (likely concurrently) to minimize total runtime.
- **FR-005**: System MUST handle 404 (Not Found) errors for repositories that no longer exist or are private.
- **FR-006**: System MUST parse standard GitHub URLs (e.g., `https://github.com/owner/repo`) to extract owner and repo names for API calls.
- **FR-007**: System MUST persist `RepoStats` data into the existing SQLite database, updating relevant columns and overwriting any existing values for those fields.
- **FR-008**: System MUST query the SQLite database to identify rows containing GitHub URLs that require enrichment (e.g., where stats are null).
- **FR-009**: System MUST maintain an `enrichment_status` column in the SQLite database for each enriched URL, indicating its processing state (e.g., `SUCCESS`, `NOT_FOUND`, `API_ERROR`, `PENDING`).

### Key Entities

- **GitHubConfig**: Stores Auth Token.
- **RepoStats**: Enriched data model (Stars, Forks, LastUpdated, Description, Language). This data will be persisted in the existing SQLite database.
- **EnrichmentResult**: Composite object linking the original Input URL to the `RepoStats` (or Error).

## Success Criteria *(mandatory)*

### Measurable Outcomes



- **SC-001**: System correctly fetches metadata for 100% of valid, public repositories in a test set.

- **SC-002**: System processes 50 repositories in under 10 seconds (assuming standard network latency and concurrent requests).

- **SC-003**: System accurately reports "Rate Limit Exceeded", processes partially, and exits cleanly when the quota is exhausted (verifiable via mock).

- **SC-004**: Unauthenticated requests work for small batches (<60) without error.



## Clarifications







### Session 2025-12-04







- Q: How should the enriched GitHub repository data be persisted or presented to the user? → A: Update SQLite



- Q: How should the system handle existing repository data during persistence? → A: Update existing entry (Overwrite)



- Q: What is the primary source of input URLs for the enrichment process? → A: From Database (Scan & Enrich)



- Q: How should the system respond when GitHub API rate limits are exceeded? → A: Exit/Fail Fast (with progress saved)



- Q: How should the enrichment status (e.g., successful, not found, API error) for each GitHub URL be tracked in the SQLite database? → A: Status Column (Add `enrichment_status` column)
















