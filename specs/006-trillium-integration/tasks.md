# Tasks: Trillium Notes Integration

**Feature**: `006-trillium-integration`
**Spec**: [specs/006-trillium-integration/spec.md](spec.md)

## Implementation Strategy
- **Phase 1: Setup**: Update configuration and interfaces.
- **Phase 2: Trillium Client**: Implement the API client.
- **Phase 3: Trillium Sink**: Implement the Sink adapter and Markdown formatting.
- **Phase 4: Configuration (US2)**: Update the setup wizard with opt-in logic.
- **Phase 5: Integration (US1)**: Wire up the CLI command.
- **Phase 6: Polish**: Verification.

---

## Phase 1: Setup

- [x] T001 [P] Update `Config` struct in `internal/config/config.go` with Trillium fields (`TrilliumURL`, `TrilliumToken`)
- [x] T002 [P] Update `ConfigLoader` in `internal/config/loader.go` to parse new fields

---

## Phase 2: Trillium Client (User Story 1)

**Goal**: Low-level API communication.

- [x] T003 [P] Implement `TrilliumClient` in `internal/adapter/trillium/client.go` (`CreateNote` method)
- [x] T004 [P] Add unit tests for `TrilliumClient` using `httptest` in `internal/adapter/trillium/client_test.go`

---

## Phase 3: Trillium Sink (User Story 1)

**Goal**: Format data and adapt to Sink interface.

- [x] T005 [P] Implement `MarkdownFormatter` helper in `internal/ui/markdown.go` (Table generation)
- [x] T006 [P] Add unit tests for `MarkdownFormatter` in `internal/ui/markdown_test.go`
- [x] T007 [US1] Implement `TrilliumSink` in `internal/adapter/trillium/sink.go` implementing `domain.Sink`

---

## Phase 4: Configuration Wizard (User Story 2)

**Goal**: Opt-in setup flow.
**Priority**: P2

- [x] T008 [US2] Update `Prompt` helper in `internal/ui/prompt.go` to support Yes/No confirmation (`AskConfirm`)
- [x] T009 [US2] Update `runSetup` in `cmd/extractor/main.go` to ask about Trillium and prompt for credentials if confirmed

---

## Phase 5: Integration (User Story 1)

**Goal**: CLI Wiring.
**Priority**: P1

- [x] T010 [US1] Update `rank` command in `cmd/extractor/main.go` to add `--sink-trillium` flag
- [x] T011 [US1] Update `runRank` logic in `cmd/extractor/main.go` to instantiate `TrilliumSink` if flag is set

---

## Phase 6: Polish & Cross-Cutting

- [x] T012 Verify error message when `--sink-trillium` is used without config
- [x] T013 Verify Markdown table rendering with special characters

---

## Dependencies

1. **T001-T002 (Config)** must complete before **T009 (Setup)** and **T011 (Rank Integration)**.
2. **T003-T004 (Client)** must complete before **T007 (Sink)**.
3. **T005-T006 (Markdown)** must complete before **T007 (Sink)**.
4. **T007 (Sink)** must complete before **T011 (Integration)**.

## Parallel Execution Examples

- **Team A**: Trillium Client & Sink (T003-T007).
- **Team B**: Configuration & Setup Wizard (T001-T002, T008-T009).
- **Team C**: Markdown Formatting (T005-T006).
