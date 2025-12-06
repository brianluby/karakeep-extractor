---
description: "Task list for Optional TUI Status"
---

# Tasks: Optional TUI Status

**Input**: Design documents from `/specs/009-optional-tui-status/`
**Prerequisites**: plan.md, spec.md, data-model.md, contracts/interfaces.md, quickstart.md

**Tests**: Tests are included for core logic adapters to support TDD as per Constitution.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2)
- Include exact file paths in descriptions

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [x] T001 Install Bubble Tea dependencies in go.mod
- [x] T002 Create directory structure for internal/ui/tui

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T003 Define ProgressReporter interface in internal/core/domain/interfaces.go
- [x] T004 Implement TextReporter adapter in internal/adapter/reporter/text_reporter.go
- [x] T005 Update Enricher to use ProgressReporter in internal/core/service/enricher.go
- [x] T006 Update Extractor to use ProgressReporter in internal/core/service/extractor.go
- [x] T007 Update main.go to use TextReporter by default in cmd/extractor/main.go
- [x] T008 Verify regression by running existing tests

**Checkpoint**: Foundation ready - Core logic decoupled from logging, legacy behavior preserved.

---

## Phase 3: User Story 1 - TUI for Enrichment Process (Priority: P1) 🎯 MVP

**Goal**: Visual progress bar for the long-running enrichment process.

**Independent Test**: Run `karakeep-extractor enrich --tui` and verify progress bar and logs.

### Tests for User Story 1

- [x] T009 [US1] Create unit tests for BubbleTeaReporter in internal/ui/tui/reporter_test.go

### Implementation for User Story 1

- [x] T010 [US1] Create basic TUI model skeleton in internal/ui/tui/model.go
- [x] T011 [US1] Implement BubbleTeaReporter adapter in internal/ui/tui/reporter.go
- [x] T012 [P] [US1] Implement Enrichment View Model (ProgressBar) in internal/ui/tui/enrich_model.go
- [x] T013 [P] [US1] Implement Log View Model (Static Tail) in internal/ui/tui/log_model.go
- [x] T014 [US1] Wire up Update loop for Enrichment events in internal/ui/tui/model.go
- [x] T015 [US1] Integrate --tui flag for enrich command in cmd/extractor/main.go

**Checkpoint**: User Story 1 functional. `enrich --tui` works.

---

## Phase 4: User Story 2 - TUI for Extraction Process (Priority: P2)

**Goal**: Visual spinner for the extraction process.

**Independent Test**: Run `karakeep-extractor extract --tui` and verify spinner and counter.

### Implementation for User Story 2

- [x] T016 [P] [US2] Implement Extraction View Model (Spinner) in internal/ui/tui/extract_model.go
- [x] T017 [US2] Update Root Model to support Extraction mode in internal/ui/tui/model.go
- [x] T018 [US2] Integrate --tui flag for extract command in cmd/extractor/main.go

**Checkpoint**: User Story 2 functional. `extract --tui` works.

---

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [x] T019 Handle terminal resize events in internal/ui/tui/model.go
- [x] T020 Ensure graceful shutdown on Ctrl+C in internal/ui/tui/program.go
- [x] T021 [P] Update documentation in docs/usage.md

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies.
- **Foundational (Phase 2)**: Depends on Setup. BLOCKS all user stories.
- **User Stories (Phase 3+)**: Depend on Foundational.
- **Polish (Phase 5)**: Depends on User Stories.

### User Story Dependencies

- **User Story 1 (P1)**: Independent after Foundational.
- **User Story 2 (P2)**: Independent after Foundational (can reuse T010/T011).

### Parallel Opportunities

- T012 and T013 can be built in parallel.
- T016 can be built in parallel with US1 tasks if T010 is ready.

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Setup + Foundational.
2. Build TUI Adapter + Enrichment View.
3. Verify `enrich --tui`.

### Incremental Delivery

1. Foundation ready (TextReporter).
2. US1: Enrichment TUI (Progress Bar).
3. US2: Extraction TUI (Spinner).