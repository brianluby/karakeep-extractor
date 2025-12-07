# Implementation Plan: Real-time TUI Statistics

**Branch**: `011-real-time-tui-stats` | **Date**: 2025-12-06 | **Spec**: [specs/011-real-time-tui-stats/spec.md](spec.md)
**Input**: Feature specification from `/specs/011-real-time-tui-stats/spec.md`

## Summary

Enhance the TUI (Text User Interface) to display real-time counters (Success, Failed, Skipped) during extraction and enrichment. This involves updating the `domain.ProgressReporter` interface to support granular status events, modifying the `BubbleTeaReporter` to dispatch new message types, and updating the TUI models (`EnrichModel`, `ExtractModel`) to track and render these statistics.

## Technical Context

**Language/Version**: Go 1.25.5
**Primary Dependencies**: `github.com/charmbracelet/bubbletea` (Existing TUI framework).
**Storage**: N/A (Transient UI state).
**Testing**: Unit tests for model updates; Manual verification for TUI layout.
**Target Platform**: CLI/TUI.
**Project Type**: Single CLI binary.
**Performance Goals**: Updates <100ms, no flicker.
**Constraints**: Must maintain existing TUI layout stability.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- [x] **CLI First**: Improves the existing TUI experience without removing CLI capabilities.
- [x] **Modular Architecture**: Updates are isolated to the UI/Reporter layer.
- [x] **Data Integrity**: Provides better visibility into data quality (failures/skips).
- [x] **TDD**: Model update logic can be unit tested.
- [x] **Simplicity**: Adds minimal complexity (just counters) for high value.

## Project Structure

### Documentation (this feature)

```text
specs/011-real-time-tui-stats/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output
│   └── interfaces.md
└── tasks.md             # Phase 2 output
```

### Source Code (repository root)

```text
src/
├── internal/
│   ├── core/
│   │   └── domain/      # Update ProgressReporter interface
│   ├── adapter/
│   │   └── reporter/    # Update TextReporter implementation
│   └── ui/
│       └── tui/         # Update BubbleTeaReporter and Models
```

**Structure Decision**: Option 1: Single project (DEFAULT). Modifying existing UI components.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| N/A       |            |                                     |