# Tasks: Karakeep AI Tags Extraction

**Feature**: Karakeep AI Tags Extraction
**Spec**: [specs/012-karakeep-ai-tags/spec.md](spec.md)
**Plan**: [specs/012-karakeep-ai-tags/plan.md](plan.md)

## Phase 1: Setup & Foundational
*Goal: Update domain models and schema to support tags.*

- [x] T001 Update `RawBookmark` struct in `internal/core/domain/entities.go` to include `Tags []string`.
- [x] T002 Update `ExtractedRepo` struct in `internal/core/domain/entities.go` to include `Tags []string`.
- [x] T003 Update `SQLiteRepository.InitSchema` in `internal/adapter/sqlite/repository.go` to create `tags` and `repo_tags` tables (including indexes).
- [x] T004 [TEST] Create integration test in `internal/adapter/sqlite/repository_test.go` to verify schema creation and migration.

## Phase 2: User Story 1 - Extract and Store Tags
*Goal: Parse tags from API and persist them.*

**Independent Test**: Run `karakeep extract` on bookmarks with tags and verify `repo_tags` table population in SQLite.

- [x] T005 [US1] [TEST] Update `client_test.go` in `internal/adapter/karakeep` to verify `FetchBookmarks` parses tags correctly (mock API response).
- [x] T006 [US1] Update `service.Extractor.Extract` in `internal/core/service/extractor.go` to pass tags from `RawBookmark` to `ExtractedRepo`.
- [x] T007 [US1] Implement `saveTags` helper in `internal/adapter/sqlite/repository.go` to normalize (lowercase/trim), insert tags, and link them to repos.
- [x] T008 [US1] Update `SQLiteRepository.Save` in `internal/adapter/sqlite/repository.go` to call `saveTags` within the transaction.
- [x] T009 [US1] [TEST] Update `extractor_test.go` (or `repository_test.go`) to verify tags are saved and linked correctly.

## Phase 3: User Story 2 - Filter by Extracted Tags
*Goal: Enable filtering by tags in the rank command.*

**Independent Test**: Run `karakeep rank --tag "golang"` and verify output.

- [x] T010 [US2] Update `SQLiteRepository.GetRankedRepos` in `internal/adapter/sqlite/repository.go` to implement filtering logic using `JOIN tags`.
- [x] T011 [US2] [TEST] Create unit test for `GetRankedRepos` with tag filter in `internal/adapter/sqlite/repository_filter_test.go`.
- [x] T012 [US2] Verify `cmd/extractor/main.go` correctly passes the existing `--tag` flag to the updated repository method (already exists, just verify plumbing).

## Phase 4: Polish & Cross-Cutting
*Goal: Ensure robustness and data quality.*

- [x] T013 Add logic to `saveTags` to handle duplicate tags gracefully (using `INSERT OR IGNORE` or similar).
- [x] T014 Verify `karakeep rank` output (table/json) includes tags if appropriate (update `internal/ui/table.go` or `export.go`).

## Implementation Strategy
- **MVP**: Complete Phase 1 and 2 first to ensure data capture. Phase 3 enables usage.
- **Incremental**: Database schema changes are additive and safe.

## Dependencies
- **Phase 1** blocks everything.
- **Phase 2** blocks **Phase 3** (need data to filter).

## Parallel Execution Examples
- **US1 & US2**: Once Phase 1 (Schema) is done, Developer A can work on `Save` logic (T006-T008) while Developer B works on `GetRankedRepos` query (T010).