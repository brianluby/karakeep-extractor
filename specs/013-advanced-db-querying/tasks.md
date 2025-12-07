# Tasks: Advanced Local Database Querying

**Feature**: Advanced Local Database Querying
**Spec**: [specs/013-advanced-db-querying/spec.md](spec.md)
**Plan**: [specs/013-advanced-db-querying/plan.md](plan.md)

## Phase 1: Setup & Foundational
*Goal: Establish data structures and interface for querying.*

- [ ] T001 Create `QueryFilter` struct in `internal/core/domain/entities.go` (including fields for ranges, dates, lang, search).
- [ ] T002 Update `RepoRepository` interface in `internal/core/domain/interfaces.go` to include `QueryRepos(ctx context.Context, filter QueryFilter) ([]ExtractedRepo, error)`.
- [ ] T003 Implement `QueryRepos` stub in `internal/adapter/sqlite/repository.go`.

## Phase 2: User Story 1 - Filter by Numeric Ranges
*Goal: Implement star and fork range filtering.*

**Independent Test**: Run `karakeep query --stars ">1000"` and verify DB query logic.

- [ ] T004 [US1] Implement SQL WHERE clause generation for numeric ranges (Stars, Forks) in `internal/adapter/sqlite/query_builder.go` (new file) or helper method.
- [ ] T005 [US1] [TEST] Create unit tests for numeric range SQL generation in `internal/adapter/sqlite/query_builder_test.go`.
- [ ] T006 [US1] Implement `parseNumericFilter` helper in `internal/ui/filter.go` to parse strings like ">1000", "10..20" into `min/max` integers.
- [ ] T007 [US1] [TEST] Create unit tests for `parseNumericFilter` in `internal/ui/filter_test.go`.
- [ ] T008 [US1] Update `QueryRepos` in `internal/adapter/sqlite/repository.go` to use the query builder.

## Phase 3: User Story 2 - Filter by Date
*Goal: Implement date-based filtering.*

**Independent Test**: Run `karakeep query --after "2024-01-01"` and verify results.

- [ ] T009 [US2] Implement SQL WHERE clause generation for Date ranges (CreatedAfter, CreatedBefore) in query builder.
- [ ] T010 [US2] [TEST] Add date range tests to `internal/adapter/sqlite/query_builder_test.go`.
- [ ] T011 [US2] Implement `parseDateFilter` helper in `internal/ui/filter.go` to parse "YYYY-MM-DD".

## Phase 4: User Story 3 - Complex Combinations & CLI Wiring
*Goal: Wire everything up to the CLI.*

**Independent Test**: Run complex `karakeep query` command.

- [ ] T012 [US3] Add `query` subcommand to `cmd/extractor/main.go` with flags (`--stars`, `--forks`, `--after`, `--before`, `--lang`, `--search`).
- [ ] T013 [US3] Implement `runQuery` function in `cmd/extractor/main.go` to parse flags into `QueryFilter` and call `repo.QueryRepos`.
- [ ] T014 [US3] Implement results display (reuse `ui.Table` or `ui.Export` logic) in `runQuery`.

## Phase 5: Polish & Cross-Cutting
*Goal: Ensure robustness.*

- [ ] T015 Implement language and search term filtering in query builder (simple exact match / LIKE).
- [ ] T016 Add validation for incompatible flags or invalid formats in `runQuery` (fail fast with helpful error).

## Implementation Strategy
- **MVP**: Numeric filters + CLI command first. Then add Dates.
- **Incremental**: Build the query builder helper progressively.

## Dependencies
- **Phase 1** blocks everything.
- **Phase 2** and **Phase 3** can be parallelized (logic is separate in builder).

## Parallel Execution Examples
- **US1 & US2**: Developer A works on Numeric Range logic (T004-T007), Developer B works on Date Range logic (T009-T011).