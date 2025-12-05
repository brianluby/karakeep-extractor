# Tasks: GitHub Repository Enrichment

**Feature**: `002-github-enrichment`
**Spec**: [specs/002-github-enrichment/spec.md](spec.md)

## Implementation Strategy
- **Phase 1: Setup**: Initialize project structure and dependencies (Go standard library).
- **Phase 2: Foundation**: Update entities, configuration, and database schema (SQLite) to support enrichment data.
- **Phase 3: Connect & Authenticate (US1)**: Implement the GitHub API client with authentication and rate limit handling.
- **Phase 4: Enrich & Persist (US2)**: Implement the core enrichment logic (worker pool), database retrieval/persistence, and CLI integration.
- **Phase 5: Polish**: Verification and cleanup.

---

## Phase 1: Setup

- [x] T001 Create GitHub adapter package structure in `internal/adapter/github`
- [x] T002 Create enrichment service package structure in `internal/core/service`

---

## Phase 2: Foundation (Blocking)

**Goal**: Prepare the data model and storage layer for enrichment data.

- [x] T003 [P] Update `RepoStats` and `EnrichmentStatus` definitions in `internal/core/domain/entities.go`
- [x] T004 [P] Add `GitHubToken` field to configuration in `internal/config/config.go`
- [x] T005 [P] Define `GitHubClient` and `Repository` interface extensions in `internal/core/port/ports.go`
- [x] T006 Implement schema migration (ALTER TABLE) for new columns in `internal/adapter/sqlite/repository.go`
- [x] T007 Implement `UpdateRepoEnrichment` method in `internal/adapter/sqlite/repository.go`
- [x] T008 Implement `GetReposForEnrichment` method in `internal/adapter/sqlite/repository.go`
- [x] T009 Add unit tests for SQLite repository methods in `internal/adapter/sqlite/repository_test.go`

---

## Phase 3: Connect and Authenticate (User Story 1)

**Goal**: reliably connect to GitHub API with/without token and handle rate limits.
**Priority**: P1

- [x] T010 [P] [US1] Implement `NewClient` factory with token support in `internal/adapter/github/client.go`
- [x] T011 [US1] Implement `GetRepoStats` with HTTP request logic in `internal/adapter/github/client.go`
- [x] T012 [US1] Implement rate limit header parsing (`X-RateLimit-*`) in `internal/adapter/github/client.go`
- [x] T013 [US1] Add `client_test.go` using `httptest` to mock GitHub API responses (Auth/Unauth/RateLimit) in `internal/adapter/github/client_test.go`

---

## Phase 4: Enrich Repository Data (User Story 2)

**Goal**: Orchestrate the fetching of data for multiple URLs and persist results.
**Priority**: P1

- [x] T014 [US2] Implement `Enricher` service struct and initialization in `internal/core/service/enricher.go`
- [x] T015 [US2] Implement worker pool logic for concurrent fetching in `internal/core/service/enricher.go`
- [x] T016 [US2] Implement "Fail Fast" logic for rate limit exceeded events in `internal/core/service/enricher.go`
- [x] T017 [US2] Integrate `GetReposForEnrichment` and `UpdateRepoEnrichment` calls in `internal/core/service/enricher.go`
- [x] T018 [US2] Add unit tests for Enricher service (mocking dependencies) in `internal/core/service/enricher_test.go`
- [x] T019 [US2] Create `enrich` CLI command definition and flags (`--token`, `--limit`, `--force`) in `cmd/extractor/main.go` (or separate command file if refactoring)
- [x] T020 [US2] Wire up Configuration, Storage, and Enricher Service in `cmd/extractor/main.go`

---

## Phase 5: Polish & Cross-Cutting

- [x] T021 Ensure 404 errors are correctly mapped to `StatusNotFound` in database
- [x] T022 Verify CLI output format matches contract (Success/Error messages)
- [x] T023 Run full integration test with seeded database and mock API

---

## Dependencies

1. **T003-T009 (Foundation)** must complete before **T017 (Enricher Persistence)**.
2. **T010-T013 (GitHub Client)** must complete before **T014 (Enricher Service)**.
3. **T014-T018 (Enricher Logic)** must complete before **T019-T020 (CLI Wiring)**.

## Parallel Execution Examples

- **Team A**: Implement SQLite Repository updates (T006-T009).
- **Team B**: Implement GitHub Client (T010-T013).
- **Team C**: Update Domain Entities & Config (T003-T004).
