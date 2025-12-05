# Feature Specification: GitHub Repository Ranking

**Feature Branch**: `003-github-ranking`
**Created**: 2025-12-04
**Status**: Draft
**Input**: User description: "Add basic github repository ranking"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - View Top Repositories (Priority: P1)

As a user, I want to see a list of my extracted GitHub repositories ordered by their popularity (stars), so I can quickly identify the most valuable resources I've saved.

**Why this priority**: This turns the raw data (enrichment) into actionable insight (ranking), fulfilling the core value proposition of "Curating" content.

**Independent Test**: Can be tested by running the `rank` command against a populated database and verifying the output is sorted by star count descending.

**Acceptance Scenarios**:

1. **Given** a database with enriched repositories, **When** I run the rank command, **Then** I see a table of repositories sorted by Stargazers count (highest to lowest).
2. **Given** a database with unenriched repositories (no stars), **When** I run the rank command, **Then** they should be omitted from the output to keep the list focused on valuable data.
3. **Given** an empty database, **When** I run the rank command, **Then** it should report "No repositories found".

---

### User Story 2 - Custom Limit and Metrics (Priority: P2)

As a user, I want to control how many results I see and potentially sort by other metrics (like recency), so I can tailor the view to my current needs.

**Why this priority**: Adds flexibility for users with large datasets who only care about the "Top 10" or "Most Recent".

**Independent Test**: Run with flags like `--limit 5` or `--sort pushed` and verify output changes.

**Acceptance Scenarios**:

1. **Given** 50 repos in DB, **When** I run `rank --limit 5`, **Then** only the top 5 are displayed.
2. **Given** enriched repos, **When** I run `rank --sort forks`, **Then** the list is sorted by Fork count instead of Stars.
3. **Given** enriched repos, **When** I run `rank --sort updated`, **Then** the list is sorted by `LastPushedAt`.

---

### Edge Cases

- **Ties**: How to handle repos with same star count? (Secondary sort by Forks, then Name).
- **Missing Data**: Repos that failed enrichment (NULL stars). Treat as 0.
- **Large Terminal output**: If output exceeds terminal height, automatically pipe output to a pager (e.g., `less`) for better navigation.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST provide a `rank` CLI command.
- **FR-002**: System MUST query the local SQLite database for enriched repository data.
- **FR-003**: System MUST default to sorting repositories by `Stars` (descending).
- **FR-004**: System MUST support a `--limit` flag (default 20) to restrict output size.
- **FR-005**: System MUST support a `--sort` flag with allowed values: `stars` (default), `forks`, `updated` (pushed_at).
- **FR-006**: System MUST display output in a formatted ASCII table containing: Rank #, Repository Name (`owner/repo`), Stars, Forks, Last Updated (relative, e.g. "2 days ago").
- **FR-007**: System MUST exclude repositories with missing or null enrichment data from the ranking output by default.
- **FR-008**: System MUST automatically detect terminal height and pipe output to a pager (like `less`) if the content exceeds the available vertical space.

### Key Entities

- **RankedRepoView**: A read-only projection of `ExtractedRepo` used for display, containing formatted strings (e.g., "2k" stars, "3mo ago").

## Success Criteria *(mandatory)*

### Measurable Outcomes



- **SC-001**: User can generate a "Top 10" list of repositories sorted by stars in under 1 second (assuming local DB).

- **SC-002**: Output table renders correctly on standard terminal widths (80 chars) without wrapping critical info.

- **SC-003**: Sorting logic correctly orders 100% of repositories based on the selected metric (verified via integration test).



## Clarifications







### Session 2025-12-04







- Q: How should large outputs exceeding terminal height be handled? → A: Paged output (automatically pipe to `less`)



- Q: How should repositories with missing enrichment data be handled in the ranking list? → A: Omit/Hide




