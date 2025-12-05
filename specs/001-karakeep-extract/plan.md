# Implementation Plan: [FEATURE]

**Branch**: `[###-feature-name]` | **Date**: [DATE] | **Spec**: [link]
**Input**: Feature specification from `/specs/[###-feature-name]/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Implement the `extract` CLI command to fetch bookmarks from a Karakeep instance, filter for valid GitHub repository URLs (normalizing to `owner/repo`), and persist them to a local SQLite database. Handles API authentication, pagination, and rate limiting with exponential backoff.

## Technical Context

**Language/Version**: Go 1.25.5  
**Primary Dependencies**: `net/http` (Stdlib), `github.com/mattn/go-sqlite3` (SQLite Driver)  
**Storage**: SQLite (local file)  
**Testing**: Go `testing` package (Stdlib)  
**Target Platform**: macOS (Dev), Linux (Target)  
**Project Type**: CLI Utility  
**Performance Goals**: <10s for 500 bookmarks  
**Constraints**: Minimal external dependencies, resilient to API failures  
**Scale/Scope**: ~1k-10k bookmarks typical load

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

*   **I. CLI First**: Passed. Feature exposes an `extract` command.
*   **II. Modular Architecture**: Passed. Separates Source (Karakeep) from Storage (SQLite).
*   **III. Data Integrity**: Passed. Spec includes logging for malformed URLs and retry logic.
*   **IV. TDD**: Passed. Tests required for parser and logic.
*   **V. Simplicity**: Passed. Minimal dependencies selected.

## Project Structure

### Documentation (this feature)

```text
specs/001-karakeep-extract/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── cli-contract.md
└── tasks.md
```

### Source Code (repository root)

```text
# Option 1: Single project (DEFAULT)
cmd/
└── extractor/
    └── main.go         # CLI entrypoint

internal/
├── adapter/
│   ├── karakeep/       # HTTP Client
│   └── sqlite/         # DB Repository
├── config/             # Env/Flag parsing
└── core/
    ├── domain/         # Entities
    └── service/        # Extraction Orchestration
```

**Structure Decision**: Option 1: Single project. Fits the existing Go project layout and Hexagonal Architecture.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| [e.g., 4th project] | [current need] | [why 3 projects insufficient] |
| [e.g., Repository pattern] | [specific problem] | [why direct DB access insufficient] |
