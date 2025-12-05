# Tasks: Configuration Wizard

**Feature**: `005-configuration-wizard`
**Spec**: [specs/005-configuration-wizard/spec.md](spec.md)

## Implementation Strategy
- **Phase 1: Setup**: Update configuration struct and adding YAML support.
- **Phase 2: Foundation**: Implement config file loading, saving, and precedence logic.
- **Phase 3: Interaction**: Implement interactive prompts for the wizard.
- **Phase 4: Integration (US1/US2)**: Wire up the `setup` command and automatic loading.
- **Phase 5: Polish**: Verification and error handling.

---

## Phase 1: Setup

- [x] T001 [P] Add `gopkg.in/yaml.v3` dependency to `go.mod`
- [x] T002 [P] Update `Config` struct in `internal/config/config.go` with YAML tags

---

## Phase 2: Foundation (Blocking)

**Goal**: Enable persistent configuration.

- [x] T003 [P] Implement `ConfigLoader` logic in `internal/config/loader.go` (Load from file, Env, Flag precedence)
- [x] T004 [P] Implement `SaveConfig` logic in `internal/config/loader.go` (File creation, 0600 permissions)
- [x] T005 Add unit tests for `ConfigLoader` (Precedence, Parsing) in `internal/config/loader_test.go`

---

## Phase 3: Interaction Logic (User Story 1)

**Goal**: Interactive setup wizard.
**Priority**: P1

- [x] T006 Implement `Prompt` helper in `internal/ui/prompt.go` (stdin reading)
- [x] T007 Add unit tests for `Prompt` helper in `internal/ui/prompt_test.go`

---

## Phase 4: Integration (User Story 1 & 2)

**Goal**: Wire up CLI commands.
**Priority**: P1

- [x] T008 [US1] Implement `runSetup` function in `cmd/extractor/main.go` (Orchestrate prompts + SaveConfig)
- [x] T009 [US2] Update `main.go` initialization to use `ConfigLoader` for all commands (extract, enrich, rank)
- [x] T010 [US1] Register `setup` command in `cmd/extractor/main.go` switch statement

---

## Phase 5: Polish & Cross-Cutting

- [x] T011 Verify file permissions are set to 0600 on create
- [x] T012 Verify helpful error message when config file is missing but flags are not provided

---

## Dependencies

1. **T001-T002 (Setup)** must complete before **T003-T004 (Foundation)**.
2. **T003-T004** must complete before **T008-T009 (Integration)**.
3. **T006 (Prompt)** must complete before **T008 (Setup Command)**.

## Parallel Execution Examples

- **Team A**: Implement Config Loader (T003-T005).
- **Team B**: Implement UI Prompts (T006-T007).
- **Team C**: Update Config Struct (T002).
