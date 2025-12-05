# Feature Specification: Karakeep Data Extraction

**Feature Branch**: `001-karakeep-extract`
**Created**: 2025-12-04
**Status**: Draft
**Input**: User description: "Implement Karakeep API client and data extraction logic."

## Clarifications

### Session 2025-12-04
- Q: Output format for extracted data? → A: Write to a local SQLite database.
- Q: How to handle duplicate/variant URLs? → A: Canonical Repository match (normalize to owner/repo).
- Q: Rate Limit handling? → A: Retry with exponential backoff.
- Q: How to handle malformed URLs? → A: Log error and skip bookmark.
- Q: How to provide Karakeep configuration (URL, Token)? → A: Environment variables and/or CLI flags.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Connect to Karakeep (Priority: P1)

As a user, I want to configure my Karakeep instance details so that the tool can access my bookmarks.

**Why this priority**: Without a connection, no data can be extracted.

**Independent Test**: Can be tested by providing valid and invalid credentials and verifying the connection status/error message.

**Acceptance Scenarios**:

1. **Given** valid API URL and Token, **When** the tool attempts to connect, **Then** it should succeed and be able to make a test request.
2. **Given** an invalid Token, **When** the tool attempts to connect, **Then** it should return an authentication error message.
3. **Given** an unreachable URL, **When** the tool attempts to connect, **Then** it should return a network error message.

---

### User Story 2 - Fetch and Filter Bookmarks (Priority: P1)

As a user, I want the tool to fetch all my bookmarks and identify the ones that are GitHub repositories, so I can analyze them later.

**Why this priority**: This is the core "Extraction" phase of the pipeline.

**Independent Test**: Can be tested by pointing the tool at a Karakeep instance (or mock) with known data and verifying the output list contains only the expected GitHub links.

**Acceptance Scenarios**:

1. **Given** a Karakeep instance with mixed bookmarks (articles, images, GitHub links), **When** the extraction runs, **Then** it should identify all bookmarks containing `github.com` URLs.
2. **Given** a user with many bookmarks (requiring pagination), **When** the extraction runs, **Then** it should traverse all pages to find all relevant links.
3. **Given** a bookmark where the main URL is NOT GitHub but the *content/summary* contains a GitHub link, **When** the extraction runs, **Then** it should extraction that link (if technically feasible per parsing rules, otherwise strictly main URL). *Assumption: Primary focus is main URL.*
4. **Given** no GitHub bookmarks exist, **When** extraction runs, **Then** it should report 0 found without error.

---

### Edge Cases

- **Rate Limiting**: Resolved. System implements retry logic with exponential backoff (e.g., 3 attempts).
- **Malformed URLs**: Resolved. System will log errors for invalid URLs and skip the affected bookmarks.
- **Duplicate Links**: Resolved. System normalizes to `owner/repo` and stores unique entries only.

## Requirements *(mandatory)*

### Functional Requirements

- FR-001: System MUST accept configuration for Karakeep Base URL and API Bearer Token via environment variables and/or CLI flags.
- **FR-002**: System MUST provide a CLI command (e.g., `extract` or `fetch`) to initiate the process.
- **FR-003**: System MUST authenticate with the Karakeep API using the provided token.
- **FR-004**: System MUST fetch bookmarks from Karakeep, handling pagination automatically.
- FR-005: System MUST filter fetched bookmarks to retain only those with `github.com` in the URL, normalizing them to the canonical `owner/repo` format (ignoring protocol, www, .git, and fragments).
- FR-006: System MUST output the extracted list of GitHub repository URLs (and associated Karakeep metadata like Title/Summary) to a local SQLite database for persistence.
- FR-007: System MUST handle common API errors (401 Unauthorized, 404 Not Found, 500 Server Error) with user-friendly error messages.
- FR-008: System MUST ensure only unique repositories (identified by `owner/repo`) are stored in the database, preventing duplicates.
- FR-009: System MUST implement exponential backoff retry logic for HTTP 429 (Too Many Requests) responses (max 3 retries).
- FR-010: System MUST log errors for malformed bookmark URLs and skip them without halting extraction.

### Key Entities

- **KarakeepConfig**: Stores URL and Token.
- **RawBookmark**: Represents the JSON structure returned by Karakeep API (id, url, title, content).
- **ExtractedRepo**: A simplified structure containing the Repository URL and origin details (Karakeep ID, Title).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of bookmarks with `github.com` as the main URL are identified from a sample set.
- **SC-002**: System handles pagination correctly for accounts with >100 bookmarks (default page size usually <100).
- **SC-003**: Extraction process for 500 bookmarks completes in under 10 seconds (network dependent, but logic should be fast).
- **SC-004**: Users receive a clear error message if the Karakeep token is invalid.