# Feature Specification: Karakeep AI Tags Extraction

**Feature Branch**: `012-karakeep-ai-tags`  
**Created**: 2025-12-06  
**Status**: Draft  
**Input**: User description: "Karakeep AI Tags Extraction **Description:** Extract AI-generated tags associated with bookmarks from the Karakeep API. **Requirement:** Research the Karakeep API response structure to identify where these tags are stored and map them to the local database schema. This will improve filtering and categorization capabilities."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Extract and Store Tags (Priority: P1)

As a user syncing my bookmarks, I want the application to automatically extract and store the AI-generated tags from Karakeep so that I can use them for filtering later.

**Why this priority**: This unlocks the core value of better categorization. Without storing the tags, they cannot be used.

**Independent Test**: Can be tested by running `karakeep extract` on a known bookmark with tags and verifying the tags appear in the local SQLite database.

**Acceptance Scenarios**:

1. **Given** a bookmark in Karakeep has AI tags (e.g., "golang", "cli"), **When** `karakeep extract` is run, **Then** those tags are saved to the `tags` table (or column) in the local database linked to the repository.
2. **Given** a bookmark has no tags, **When** extracted, **Then** the system handles the empty list gracefully without errors.
3. **Given** a bookmark is updated in Karakeep with new tags, **When** extracted again, **Then** the local database updates to reflect the new tag set.

---

### User Story 2 - Filter by Extracted Tags (Priority: P2)

As a user ranking repositories, I want to filter the output using the tags extracted from Karakeep (e.g., `--tag "golang"`).

**Why this priority**: Makes the extracted data actionable.

**Independent Test**: Run `karakeep rank --tag "python"` and verify only repositories with the "python" tag are displayed.

**Acceptance Scenarios**:

1. **Given** a populated database with tagged repos, **When** `karakeep rank --tag "rust"` is run, **Then** only repositories containing the "rust" tag are shown.
2. **Given** no repos match the tag, **When** the command is run, **Then** an empty list (or "no results") message is displayed.

---

### Edge Cases

- **Tag Format**: Tags might contain spaces or special characters. The system should normalize them (e.g., lowercase) for consistent filtering.
- **Duplicate Tags**: The API might return duplicates; the system should store unique tags per repo.
- **Database Migration**: Existing users will have a database without a `tags` column/table. A migration strategy is required.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST parse the `tags` (or equivalent) field from the Karakeep API response for each bookmark.
- **FR-002**: The system MUST store these tags in the local SQLite database, associated with the `ExtractedRepo`.
- **FR-003**: The system MUST normalize tags (trim whitespace, lowercase) before storage.
- **FR-004**: The `karakeep rank` command MUST support filtering by these locally stored tags via the `--tag` flag (updating the existing behavior which might only filter by title/desc or be a placeholder).
- **FR-005**: The system MUST automatically migrate the database schema to support tags if the table/column is missing.

### Key Entities *(include if feature involves data)*

- **Tag**: A string label (e.g., "devops").
- **Repository**: Updated to include a list of `Tags` or a relation to a `Tags` table.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of tags present in the Karakeep API response for a bookmark are persisted locally after extraction.
- **SC-002**: Users can filter by a tag found in Karakeep and get correct results within 2 seconds.
- **SC-003**: Database migration runs automatically and successfully on the next CLI execution for existing users.