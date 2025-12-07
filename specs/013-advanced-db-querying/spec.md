# Feature Specification: Advanced Local Database Querying

**Feature Branch**: `013-advanced-db-querying`  
**Created**: 2025-12-06  
**Status**: Draft  
**Input**: User description: "Advanced Local Database Querying **Description:** Add a powerful CLI command to query the local SQLite database with flexible filters. **Capabilities:** Support filtering by multiple keywords, star count ranges (e.g., `>1000`), fork counts, dates (e.g., `created_after=2024-01-01`), and combinations of these criteria."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Filter by Numeric Ranges (Priority: P1)

As a power user, I want to filter repositories by specific numeric ranges (stars, forks) so I can find high-quality or niche projects.

**Why this priority**: This is the core "advanced" functionality requested.

**Independent Test**: Can be tested by running `karakeep query --stars ">1000"` and verifying all returned repos have >1000 stars.

**Acceptance Scenarios**:

1. **Given** a database of repos, **When** running `karakeep query --stars ">500"`, **Then** only repos with >500 stars are shown.
2. **Given** a database of repos, **When** running `karakeep query --forks "<50"`, **Then** only repos with <50 forks are shown.
3. **Given** a database of repos, **When** running `karakeep query --stars "100..200"`, **Then** only repos with stars between 100 and 200 (inclusive) are shown.

---

### User Story 2 - Filter by Date (Priority: P2)

As a user looking for recent projects, I want to filter by date (extracted date or last pushed date).

**Why this priority**: Important for finding fresh content vs stale bookmarks.

**Independent Test**: Run `karakeep query --after "2024-01-01"` and check the `FoundAt` or `LastPushedAt` date.

**Acceptance Scenarios**:

1. **Given** a database, **When** running `karakeep query --after "2024-01-01"`, **Then** only repos found after that date are shown.
2. **Given** a database, **When** running `karakeep query --before "2023-01-01"`, **Then** only repos found before that date are shown.

---

### User Story 3 - Complex Combinations (Priority: P3)

As a user, I want to combine multiple filters to create a specific search query.

**Why this priority**: Demonstrates the power of the new command.

**Independent Test**: Run `karakeep query --stars ">1000" --lang "Go" --after "2024-01-01"`.

**Acceptance Scenarios**:

1. **Given** a database, **When** combining `--stars >1000` AND `--lang Go`, **Then** only Go repos with >1000 stars are returned.

---

### Edge Cases

- **Invalid Number Format**: `--stars "abc"` should error gracefully.
- **Invalid Date Format**: `--after "yesterday"` should error (unless we support natural language, but ISO 8601 is safer for MVP).
- **No Results**: Should display "No repositories found matching criteria."

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST provide a `karakeep query` command.
- **FR-002**: The query command MUST support `--stars` argument accepting operators `>`, `<`, `>=`, `<=`, and ranges `min..max`.
- **FR-003**: The query command MUST support `--forks` argument with the same operators as stars.
- **FR-004**: The query command MUST support `--after` and `--before` arguments accepting dates (YYYY-MM-DD) filtering on `FoundAt` (extraction date) by default.
- **FR-005**: The query command MUST support `--lang` (exact or case-insensitive match).
- **FR-006**: The query command MUST support `--search` for keyword matching in Title/Description.
- **FR-007**: All filters MUST be additive (AND logic).

### Key Entities *(include if feature involves data)*

- **QueryFilter**: A struct representing the parsed filter criteria.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Complex queries (3+ filters) execute in under 500ms on a database of 10,000 items.
- **SC-002**: Users can construct a valid query for "popular recent Go projects" without reading documentation (intuitive flags).