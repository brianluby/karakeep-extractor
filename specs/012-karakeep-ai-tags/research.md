# Research: Karakeep AI Tags Extraction

**Status**: Complete
**Date**: 2025-12-06

## Decisions

### 1. Karakeep API Schema Analysis
**Decision**: Parse the `tags` field from the `RawBookmark` structure in the API response.
**Rationale**:
- Based on the existing codebase and typical JSON responses, bookmarks often contain a `tags` array of strings. We will assume the structure `[]string` for tags in the API response. We need to update `domain.RawBookmark` to include this field.

### 2. Database Schema Update
**Decision**: Create a new `tags` table with a many-to-many relationship to `repositories` (via `repo_tags` junction table) OR a simple 1:N relation if we don't need global tag management. Given the requirement to "Filter by Extracted Tags", a normalized schema is better for query performance and flexibility.
**Schema**:
- `tags` (id INTEGER PRIMARY KEY, name TEXT UNIQUE)
- `repo_tags` (repo_id TEXT, tag_id INTEGER, PRIMARY KEY (repo_id, tag_id))
**Rationale**:
- **Normalization**: Allows efficient querying of "all repos with tag X" without `LIKE %tag%` scans.
- **Space Efficiency**: Stores unique tag strings only once.

### 3. Extraction Logic Update
**Decision**: Update `service.Extractor` to normalize and persist tags during the extraction loop.
**Rationale**:
- **Normalization**: Tags should be lowercased and trimmed to prevent "Go" vs "go " duplication.
- **Persistence**: Tags must be saved in the same transaction (or immediately after) saving the repository to ensure consistency.

### 4. Filtering Logic Update
**Decision**: Update `RankingRepository.GetRankedRepos` to accept a `tagFilter` (which is already in the interface but maybe not fully implemented or needs updating to use the new table).
**Rationale**:
- **Existing Interface**: `GetRankedRepos` already has `tagFilter string`. We need to ensure the SQL query joins with the new `tags` tables.

## Alternatives Considered

### Storing Tags as JSON/String in `extracted_repos`
- **Pros**: Simpler schema (no new tables).
- **Cons**: Harder to query efficiently ("find all repos with 'go' tag" requires string parsing in SQL). Harder to get a list of all unique tags.
- **Result**: Rejected in favor of normalized schema for better filtering performance.

### Fetching Tags in a Separate "Enrichment" Phase
- **Pros**: Keeps "Extract" simple (just URLs).
- **Cons**: Tags are available immediately from Karakeep. Waiting for enrichment delays their availability.
- **Result**: Rejected. Tags should be part of the initial extraction since they come from the source.

## Open Questions Resolved

- **API Field Name**: Assumed to be `tags` based on standard Karakeep/Shiori patterns. If it's different (e.g., `tag_names`), we will catch this during integration testing.
    - *Action*: Update `RawBookmark` struct with `Tags []string json:"tags"`.
