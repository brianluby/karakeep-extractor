# Feature Specification: Filtering by Tag

**Feature Branch**: `007-filter-by-tag`
**Created**: 2025-12-04
**Status**: Draft
**Input**: User description: "Let's implement the feature Filtering by Tag. karakeep-extractor run [--tag "dev-tools"] We extract tags (or rather, "content") from Karakeep, but we don't expose a way to filter the rank output by these tags/keywords."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Filter by Keyword/Tag (Priority: P1)

As a user, I want to filter my ranked repository list by a specific keyword or tag (derived from the Karakeep content), so I can focus on a specific topic like "dev-tools" or "database".

**Why this priority**: Increases the utility of the tool by allowing users to slice their data, rather than viewing one monolithic list.

**Independent Test**: Run `karakeep rank --tag "database"` and verify the output only contains repositories where the extracted content/title matches "database".

**Acceptance Scenarios**:

1. **Given** a database with mixed repositories, **When** I run `karakeep rank --tag "go"`, **Then** the output lists only repositories containing "go" in their title or content field.
2. **Given** no repositories match the tag, **When** I run the command, **Then** the output should display "No repositories found matching tag 'X'".
3. **Given** a tag search, **When** viewing the output, **Then** the ranking order (e.g., stars) should be preserved within the filtered subset.

---

### Edge Cases

- **Case Sensitivity**: Should "Go" match "go"? (Assumption: Case-insensitive for better UX).
- **Partial Matches**: Should "tool" match "devtools"? (Assumption: Yes, simple substring match).
- **Multiple Tags**: Can I filter by "go" AND "web"? (Scope limit: Start with single tag support first).

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST support a `--tag` (or `--filter`) flag for the `rank` command.
- **FR-002**: System MUST query the local SQLite database to filter `ExtractedRepo` entities where the `Title` or `Description` (from GitHub enrichment) contains the specified string.
- **FR-003**: The filtering MUST be case-insensitive.
- **FR-004**: The system MUST apply the filter *before* or *during* the ranking/sorting process to ensure the `--limit` applies to the *matching* set (e.g. "Top 10 GO projects", not "Top 10 projects, filtered by GO").

### Key Entities

- **ExtractedRepo**: Existing entity. Filtering applies to `Title` and `Description` (from GitHub enrichment). *Note: The original 'Content' from Karakeep is not currently persisted in the `extracted_repos` table. For this feature, filtering will target `Title` and `Description` as a fallback.*

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Filtering works instantly (< 200ms extra latency) on a database of 1000 items.
- **SC-002**: 100% of displayed results contain the filter string in Title or Description.