# Tasks: Karakeep Data Extraction

**Input**: Design documents from `/specs/001-karakeep-extract/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Tests are MANDATORY for critical logic as per Constitution Principle IV (TDD).

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2)
- Include exact file paths in descriptions

## Path Conventions

- `cmd/`: Application entry points
- `internal/`: Private application and library code
- `internal/adapter/`: Interfaces to external systems (DB, API)
- `internal/core/`: Domain logic and services
- `internal/config/`: Configuration logic

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [x] T001 Create project directory structure (cmd/, internal/...) per plan
- [x] T002 Initialize Go module `github.com/brianluby/karakeep-extractor`
- [x] T003 [P] Add `github.com/mattn/go-sqlite3` dependency
- [x] T004 [P] Create `internal/core/domain/entities.go` with empty structs (placeholders)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T005 Create `internal/config/config.go` for Env/Flag parsing (Karakeep URL, Token, DB Path)
- [x] T006 Create `internal/config/config_test.go` to verify flag/env precedence
- [x] T007 Define `KarakeepConfig` struct in `internal/core/domain/entities.go`
- [x] T008 Define `RawBookmark` struct in `internal/core/domain/entities.go`
- [x] T009 Define `ExtractedRepo` struct in `internal/core/domain/entities.go`
- [x] T010 Define `BookmarkSource` interface in `internal/core/domain/interfaces.go`
- [x] T011 Define `RepoRepository` interface in `internal/core/domain/interfaces.go`

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Connect to Karakeep (Priority: P1) 🎯 MVP

**Goal**: Configure connection and authenticate with Karakeep API.

**Independent Test**: Run CLI with valid/invalid credentials and verify connection.

### Tests for User Story 1 ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [x] T012 [P] [US1] Create `internal/adapter/karakeep/client_test.go` verifying connection logic (mock server)
- [x] T013 [P] [US1] Add test case for exponential backoff retry logic in `internal/adapter/karakeep/client_test.go`

### Implementation for User Story 1

- [x] T013 [US1] Implement `NewClient` in `internal/adapter/karakeep/client.go` using `net/http`
- [x] T014 [US1] Implement `Ping` or initial `Fetch` method in `internal/adapter/karakeep/client.go` to verify auth
- [x] T015 [US1] Create `cmd/extractor/main.go` basic structure to load config and init client
- [x] T016 [US1] Implement graceful error handling for 401/404/500 in `internal/adapter/karakeep/client.go`

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently

---

## Phase 4: User Story 2 - Fetch and Filter Bookmarks (Priority: P1)

**Goal**: Fetch bookmarks, filter for GitHub URLs, and save to SQLite.

**Independent Test**: Run extraction against mock/real API and check SQLite DB content.

### Tests for User Story 2 ⚠️

- [x] T017 [P] [US2] Create `internal/core/service/extractor_test.go` verifying filter/normalization logic
- [x] T018 [P] [US2] Create `internal/adapter/sqlite/repository_test.go` verifying DB schema/inserts

### Implementation for User Story 2

- [x] T019 [US2] Implement `FetchBookmarks` with pagination support in `internal/adapter/karakeep/client.go`
- [x] T020 [US2] Implement exponential backoff retry middleware in `internal/adapter/karakeep/client.go`
- [x] T021 [US2] Implement `SQLiteRepository` in `internal/adapter/sqlite/repository.go` (Schema creation + Save)
- [x] T022 [US2] Implement `Extract` service in `internal/core/service/extractor.go` (Orchestration: Fetch -> Filter -> Normalize -> Save)
- [x] T023 [US2] Wire up `Extract` service in `cmd/extractor/main.go`
- [x] T024 [US2] Add logging for malformed URLs in `internal/core/service/extractor.go` (skip logic)

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently

---

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [x] T025 [P] Update `README.md` with build instructions
- [x] T026 Verify `quickstart.md` steps manually
- [x] T027 Ensure error messages are user-friendly (no stack traces in CLI output)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies
- **Foundational (Phase 2)**: Depends on Setup
- **User Story 1 (Phase 3)**: Depends on Foundational
- **User Story 2 (Phase 4)**: Depends on Foundational (can technically run parallel to US1 if mocked, but US1 is logical precursor)

### User Story Dependencies

- **User Story 1 (P1)**: Needs Config and Entities (Foundational)
- **User Story 2 (P1)**: Needs Entities and Interfaces (Foundational). Integration needs US1's client.

### Parallel Opportunities

- T003, T004 (Setup)
- T012 (US1 Test) and T013 (US1 Impl) can start once Foundational is done
- T017, T018 (US2 Tests) can start once Foundational is done

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1 & 2.
2. Implement US1 (Connect).
3. Verify connection succeeds.

### Incremental Delivery

1. Complete US1.
2. Implement US2 (Extraction logic).
3. Verify full flow: Config -> Connect -> Fetch -> Filter -> Save.
