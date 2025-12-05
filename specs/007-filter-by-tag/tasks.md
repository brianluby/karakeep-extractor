# Tasks: Filtering by Tag

**Feature**: `007-filter-by-tag`
**Spec**: [specs/007-filter-by-tag/spec.md](spec.md)

## Implementation Strategy
- **Phase 1: Setup**: Update interfaces.
- **Phase 2: Foundation**: Implement SQL filtering logic.
- **Phase 3: Integration (US1)**: Wire up CLI command and service logic.
- **Phase 4: Polish**: Verification.

---

## Phase 1: Setup

- [x] T001 [P] Update `RankingRepository` interface in `internal/core/domain/interfaces.go` to accept `filterTag` string

---

## Phase 2: Foundation (Blocking)

**Goal**: Enable database to filter results efficiently.

- [x] T002 [P] Update `GetRankedRepos` in `internal/adapter/sqlite/repository.go` to include `LIKE` clause for title/description
- [x] T003 [P] Add unit tests for filtered queries in `internal/adapter/sqlite/repository_test.go`

---

## Phase 3: Integration (User Story 1)

**Goal**: Expose filtering via CLI.
**Priority**: P1

- [x] T004 [US1] Update `Ranker.Rank` method signature in `internal/core/service/ranker.go` to accept filter
- [x] T005 [US1] Update `Rank` implementation to pass filter to repository
- [x] T006 [US1] Update `rank` CLI command in `cmd/extractor/main.go` to parse `--tag` flag
- [x] T007 [US1] Update `runRank` in `cmd/extractor/main.go` to pass tag to service

---

## Phase 4: Polish & Cross-Cutting

- [x] T008 Verify "No repositories found matching tag" behavior (empty result set handles this implicitly)
- [x] T009 Verify case-insensitivity manually or via test

---

## Dependencies

1. **T001 (Interface)** must complete before **T002 (Implementation)** and **T004 (Service)**.
2. **T002-T003** must complete before **T005 (Service Logic)**.
3. **T004-T005** must complete before **T007 (CLI Wiring)**.

## Parallel Execution Examples

- **Team A**: SQL Implementation (T002-T003).
- **Team B**: Service & CLI Updates (T004-T007) - *Blocked by T001*.
