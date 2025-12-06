# Implementation Plan - Optional TUI Status

**Feature**: Optional TUI Status
**Status**: Approved
**Feature Branch**: `009-optional-tui-status`

## Technical Context

**Language:** Go 1.25.5
**New Dependencies:**
- `github.com/charmbracelet/bubbletea` (Main TUI loop)
- `github.com/charmbracelet/bubbles` (Components: progress, spinner)
- `github.com/charmbracelet/lipgloss` (Styling)

**Core Architecture:**
- **Observer Pattern:** A new `domain.ProgressReporter` interface will decouple `core/service` from `ui`.
- **UI-Driven Loop:** `main.go` will instantiate the TUI, which then spawns the worker routine.
- **Legacy Support:** A `TextReporter` implementation will maintain existing behavior when `--tui` is missing.

## Constitution Check

- [x] **CLI First:** TUI is an *option* (`--tui`), not the default. Core logic remains CLI-driven.
- [x] **Modular Architecture:** Logic is decoupled via `ProgressReporter`. UI code resides in `internal/ui/tui`.
- [x] **Data Integrity:** TUI visualization does not alter data processing logic.
- [x] **TDD:** New Reporters and TUI models will be unit tested.

## Phased Implementation

### Phase 1: Core Abstractions & Refactoring
*Goal: Decouple logging/progress from services.*

1.  **Define Interface:** Create `domain.ProgressReporter`.
2.  **Implement No-Op/Text Reporter:** Create `internal/adapter/reporter/text_reporter.go` that mimics current `log.Printf` behavior.
3.  **Refactor Services:** Update `Enricher.EnrichBatch` and `Extractor.Extract` to accept `ProgressReporter` and replace direct `log.Printf` calls with `reporter.Log()` or `reporter.Increment()`.
4.  **Update Main:** Update `cmd/extractor/main.go` to pass the `TextReporter` by default.
5.  **Verify:** Run existing tests to ensure no regression in logic.

### Phase 2: TUI Implementation
*Goal: Build the Bubble Tea models.*

1.  **Scaffold TUI Package:** Create `internal/ui/tui/model.go`.
2.  **Implement BubbleTeaReporter:** An adapter that converts `ProgressReporter` calls into `tea.Msg` events sent to the program.
3.  **Build Models:**
    - `EnrichModel`: With `bubbles/progress`.
    - `ExtractModel`: With `bubbles/spinner`.
    - `RootModel`: Switches between them.
4.  **Implement View:** Use `lipgloss` to style the output (bars, logs).
5.  **Integration:** Wire up `cmd/extractor/main.go` to instantiate `tea.Program` when `--tui` is set.

### Phase 3: Polish & Testing
*Goal: Ensure stability and edge case handling.*

1.  **Resizing:** Test terminal resize events.
2.  **Errors:** Verify `MsgFatal` handles graceful shutdown and restores terminal.
3.  **CI/CD Check:** Ensure `--tui` is not default (Constitution check).

## Validated Assets
- [x] `research.md`
- [x] `data-model.md`
- [x] `contracts/interfaces.md`
- [x] `quickstart.md`
