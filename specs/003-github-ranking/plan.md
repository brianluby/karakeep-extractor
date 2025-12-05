# Implementation Plan: GitHub Repository Ranking

**Branch**: `003-github-ranking` | **Date**: 2025-12-04 | **Spec**: [specs/003-github-ranking/spec.md](spec.md)
**Input**: Feature specification from `/specs/003-github-ranking/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

This feature implements the `rank` CLI command to display extracted GitHub repositories sorted by popularity (stars), forks, or last update time. It includes a flexible sorting and limiting mechanism, utilizes the existing SQLite database, and introduces a new presentation layer for formatted table output.

## Technical Context

**Language/Version**: Go 1.25.5
**Primary Dependencies**: `net/http` (Stdlib), `encoding/json` (Stdlib), `github.com/mattn/go-sqlite3` (Existing)
**Storage**: SQLite (Read-only access for this feature)
**Testing**: Go `testing` package
**Target Platform**: CLI (Cross-platform: macOS, Linux, Windows)
**Project Type**: CLI Tool
**Performance Goals**: Generate "Top 10" list < 1s
**Constraints**: Output formatting must handle standard terminal widths; Paged output for large lists.
**Scale/Scope**: Local database, typical 100-5000 rows.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

*   **CLI First**: ✅ Yes, adds `rank` command to CLI.
*   **Modular Architecture**: ✅ Yes, introduces `internal/ui` or similar for presentation, keeping logic in `internal/core/service`.
*   **Data Integrity**: ✅ Yes, Read-only operation on DB. Handles NULLs/missing data gracefully (as 0).
*   **TDD**: ✅ Yes, sorting logic will be unit tested.
*   **Simplicity**: ✅ Yes, using standard library or minimal helpers. avoiding complex TUI frameworks if possible, or using established one if needed for table rendering.

## Project Structure

### Documentation (this feature)

```text
specs/003-github-ranking/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
cmd/
└── extractor/
    └── main.go            # Wiring up the new Ranker service and CLI command

internal/
├── adapter/
│   └── sqlite/
│       ├── repository.go  # UPDATE: Add GetRankedRepos method
│       └── repository_test.go
├── core/
│   ├── domain/
│   │   └── entities.go    # No major changes expected, maybe RepoStats structs reuse
│   ├── port/              # UPDATE: Add RankingRepository interface
│   └── service/
│       ├── ranker.go      # NEW: Business logic for sorting/filtering (if not done in SQL)
│       └── ranker_test.go
├── ui/                    # NEW: Presentation layer
│   ├── table.go           # Table rendering logic
│   └── table_test.go
```

**Structure Decision**: Standard Go Project Layout (cmd/internal). Introducing `internal/ui` for output formatting to separate presentation from logic.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| None | | |
