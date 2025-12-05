# Tasks: GitHub Repository Ranking

**Feature**: `003-github-ranking`
**Spec**: [specs/003-github-ranking/spec.md](spec.md)

## Implementation Strategy
- **Phase 1: Setup**: Update interfaces and package structure.
- **Phase 2: Foundation**: Implement SQLite query logic for sorting and limiting (SQL).
- **Phase 3: Presentation**: Implement the table rendering and paging logic.
- **Phase 4: Ranking (US1/US2)**: Wire up the CLI command and service logic.
- **Phase 5: Polish**: Verification and cleanup.

---

## Phase 1: Setup

- [x] T001 [P] Update `RankingRepository` interface in `internal/core/domain/interfaces.go` (Move from contracts to domain package)
- [x] T002 [P] Create `internal/ui` package structure for presentation logic

---

## Phase 2: Foundation (Blocking)

**Goal**: Enable database to return sorted and limited results efficiently.

- [x] T003 [P] Implement `GetRankedRepos` in `internal/adapter/sqlite/repository.go` with dynamic query construction for sorting
- [x] T004 [P] Add unit tests for `GetRankedRepos` (sorting/limiting) in `internal/adapter/sqlite/repository_test.go`

---

## Phase 3: Presentation Logic

**Goal**: Format data into a clean table and handle pagination.

- [x] T005 Implement `TableRenderer` in `internal/ui/table.go` using `text/tabwriter`
- [x] T006 Implement TTY detection and paging logic (piping to `less`) in `internal/ui/pager.go`
- [x] T007 Add unit tests for `TableRenderer` (buffer output check) in `internal/ui/table_test.go`

---

## Phase 4: Ranking Logic (User Story 1 & 2)

**Goal**: Orchestrate the ranking command.
**Priority**: P1 & P2

- [x] T008 [US1] Implement `Ranker` service in `internal/core/service/ranker.go` to bridge Repo and UI
- [x] T009 [US1] Add unit tests for `Ranker` service in `internal/core/service/ranker_test.go`
- [x] T010 [US1] Create `rank` CLI command definition and flags (`--limit`, `--sort`) in `cmd/extractor/main.go`
- [x] T011 [US2] Wire up CLI flags to service call in `cmd/extractor/main.go`

---

## Phase 5: Polish & Cross-Cutting

- [x] T012 Verify "No repositories found" empty state
- [x] T013 Verify relative time formatting (e.g. "2d ago") in table output
- [x] T014 Verify behavior when piping output to file (paging should be disabled)

---

## Dependencies

1. **T003-T004 (Foundation)** must complete before **T008 (Ranker Service)**.
2. **T005-T007 (Presentation)** must complete before **T008 (Ranker Service)**.
3. **T008 (Ranker)** must complete before **T010 (CLI)**.

## Parallel Execution Examples

- **Team A**: Implement SQLite Repository updates (T003-T004).
- **Team B**: Implement UI Table Renderer & Pager (T005-T007).
- **Team C**: CLI Command setup (T010).
