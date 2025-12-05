# Research & Decision Log

## 1. Filtering Mechanism

**Context**: We need to filter repositories by a "tag" or keyword.
**Question**: Should we add a new `tags` table (many-to-many) or fuzzy search text?

**Decision**: **Fuzzy Text Search (`LIKE %tag%`)**

**Rationale**:
*   **Data Availability**: We currently extract "Title" and potentially "Content" (summary) from Karakeep. We don't have structured tags from the source yet (or rather, the "content" *is* the tag source).
*   **Simplicity**: SQL `LIKE` query is sufficient for the scale (100-5000 rows) and requires no schema migration if we search existing columns.
*   **User Intent**: Users want to find "python" things. Searching title/content covers this broadly.

## 2. Where to Filter

**Context**: Filter in SQL vs Go Application?
**Question**: Should we fetch all "SUCCESS" repos and filter in Go, or filter in SQL?

**Decision**: **SQL**

**Rationale**:
*   **Performance**: SQLite is optimized for this.
*   **Pagination/Limit**: If we filter in Go, the `--limit` flag logic becomes complex (fetching 1000 to find 10 matches). Filtering in SQL ensures `LIMIT 10` returns 10 matches (if they exist).

**Alternatives Considered**:
*   Filter in Go: Easier to do complex regex matching, but less efficient for limiting results.

## 3. Content Persistence Check

**Context**: Does the `extracted_repos` table have a column for the raw content/summary to search against?
**Investigation**: Checked `001` specs. `RawBookmark` has `Content`. `ExtractedRepo` entity has `Title`, `URL`, `SourceID`.
**Gap**: It seems `Content` (the summary/description from Karakeep) was NOT persisted to `extracted_repos` in feature 001. `RepoStats` adds `description` (from GitHub).
**Decision**:
    *   **Phase A (Current)**: Filter against `Title` (from Karakeep) and `Description` (from GitHub enrichment).
    *   **Phase B (Future)**: Add `content` column to `extracted_repos` and backfill if needed.
    *   *Refinement*: For this feature, we will search: `title` OR `description` (GitHub). This covers most "tagging" needs (e.g. "python", "cli").

**Update**: The user request specifically mentions "We extract tags (or rather, 'content') from Karakeep". If this content isn't in the DB, we can't search it.
*   **Action**: I will check the codebase (via tasks) to confirm if `content` or `summary` column exists.
*   *Assumption*: If it doesn't exist, searching GitHub `Description` is a valid fallback and high value.
