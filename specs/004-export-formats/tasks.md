# Tasks: Export Formats & API Sink

**Feature**: `004-export-formats`
**Spec**: [specs/004-export-formats/spec.md](spec.md)

## Implementation Strategy
- **Phase 1: Setup**: Update interfaces.
- **Phase 2: Export Logic**: Implement JSON and CSV formatting.
- **Phase 3: Sink Logic**: Implement generic HTTP sink adapter.
- **Phase 4: Integration (US1/US2)**: Update `Ranker` service and CLI command.
- **Phase 5: Polish**: Verification.

---

## Phase 1: Setup

- [x] T001 [P] Update `Exporter` and `Sink` interfaces in `internal/core/domain/interfaces.go` (Move from contracts/spec to domain package)

---

## Phase 2: Export Logic (User Story 1)

**Goal**: Support JSON/CSV output.
**Priority**: P1

- [x] T002 [P] [US1] Implement `JSONExporter` in `internal/ui/export.go`
- [x] T003 [P] [US1] Implement `CSVExporter` in `internal/ui/export.go`
- [x] T004 [P] [US1] Implement `ExporterFactory` or simple switch logic in `internal/ui/export.go` to select format
- [x] T005 [US1] Add unit tests for exporters in `internal/ui/export_test.go`

---

## Phase 3: Sink Logic (User Story 2)

**Goal**: Support generic API POST sink.
**Priority**: P2

- [x] T006 [P] [US2] Implement `HTTPSink` in `internal/adapter/http/sink.go` to handle URL and custom headers.
- [x] T007 [US2] Add unit tests for `HTTPSink` using `httptest` in `internal/adapter/http/sink_test.go`, covering header application and error handling.

---

## Phase 4: Integration

**Goal**: Wire up new capabilities to the `rank` command.

- [x] T008 Update `Ranker` service in `internal/core/service/ranker.go` to accept optional `Sink` and `Exporter` (Strategy pattern).
- [x] T009 Update `Ranker` service tests in `internal/core/service/ranker_test.go` to cover export and sink logic.
- [x] T010 Update `rank` CLI command in `cmd/extractor/main.go` to parse `--format`, `--sink-url`, and repeatable `--sink-header` flags
- [x] T011 Wire up `Exporter` selection and `HTTPSink` initialization in `cmd/extractor/main.go`

---

## Phase 5: Polish & Cross-Cutting

- [x] T012 Verify empty list behavior for JSON (empty array `[]`) and CSV (header only)
- [x] T013 Verify large payload handling (basic test with ~50 items)

---

## Dependencies

1. **T001 (Interfaces)** must complete before **T002-T007 (Implementations)**.
2. **T002-T007** must complete before **T008 (Service Integration)**.
3. **T008** must complete before **T010-T011 (CLI Integration)**.

## Parallel Execution Examples

- **Team A**: Implement Exporters (T002-T005).
- **Team B**: Implement Sink Adapter (T006-T007).
- **Team C**: Update Domain Interfaces (T001).
