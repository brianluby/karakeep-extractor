# Tasks: Real-time TUI Statistics

**Feature**: Real-time TUI Statistics
**Spec**: [specs/011-real-time-tui-stats/spec.md](spec.md)
**Plan**: [specs/011-real-time-tui-stats/plan.md](plan.md)

## Phase 1: Setup & Foundational
*Goal: Update core interfaces and data structures to support statistics tracking.*

- [x] T001 Update `domain.ProgressReporter` interface in `internal/core/domain/interfaces.go` to include `RecordSuccess`, `RecordFailure`, and `RecordSkipped` methods.
- [x] T002 Update `TextReporter` implementation in `internal/adapter/reporter/text_reporter.go` to satisfy the new interface (stub or simple log).
- [x] T003 Define `MsgSuccess`, `MsgFailure`, and `MsgSkipped` types in `internal/ui/tui/model.go`.
- [x] T004 Update `BubbleTeaReporter` in `internal/ui/tui/reporter.go` to dispatch the new Bubble Tea messages.

## Phase 2: User Story 1 - Real-time Enrichment Feedback
*Goal: Display real-time counters during enrichment process.*

**Independent Test**: Run `karakeep enrich --tui` and observe Success/Failed/Skipped counters.

- [x] T005 [US1] Create `ProgressStats` struct in `internal/ui/tui/model.go` (or shared file) to track counts.
- [x] T006 [US1] [TEST] Create unit tests in `internal/ui/tui/enrich_model_test.go` to verify `Update` handles `MsgSuccess`, `MsgFailure`, `MsgSkipped` correctly.
- [x] T007 [US1] Update `EnrichModel` in `internal/ui/tui/enrich_model.go` to embed `ProgressStats` and handle new update messages.
- [x] T008 [US1] Update `EnrichModel.View` in `internal/ui/tui/enrich_model.go` to render the statistics (Success, Failed, Skipped).
- [x] T009 [US1] Update `service.Enricher` (and related calls) to call `reporter.RecordSuccess/Failure/Skipped` instead of just logging strings (requires identifying where these events happen in `internal/core/service/enricher.go`).

## Phase 3: User Story 2 - Real-time Extraction Feedback
*Goal: Display real-time counters during extraction process.*

**Independent Test**: Run `karakeep extract --tui` and observe processed counters.

- [x] T010 [US2] Update `ExtractModel` in `internal/ui/tui/extract_model.go` to embed `ProgressStats` and handle update messages.
- [x] T011 [US2] Update `ExtractModel.View` in `internal/ui/tui/extract_model.go` to render the statistics.
- [x] T012 [US2] Update `service.Extractor` in `internal/core/service/extractor.go` to call `reporter.RecordSuccess/Failure/Skipped`.

## Phase 4: Polish & Cross-Cutting
*Goal: Ensure UI stability and clean up.*

- [x] T013 Verify 4-digit number layout stability in `EnrichModel` and `ExtractModel` views (manual check or unit test layout).
- [x] T014 Update `TextReporter` to provide meaningful feedback (e.g. dots or summary) for the new methods if appropriate.

## Implementation Strategy
- **MVP**: Complete Phase 1 and Phase 2. This covers the most critical user story (Enrichment).
- **Incremental**: Phase 3 (Extraction) reuses the same UI patterns established in Phase 2.

## Dependencies
- **Phase 1** is blocking for everything.
- **Phase 2** and **Phase 3** are independent but share the underlying `ProgressStats` structure.

## Parallel Execution Examples
- **US1 & US2 UI**: Once Phase 1 is done, Developer A can work on `EnrichModel` (T007, T008) while Developer B works on `ExtractModel` (T010, T011).