# Implementation Plan: Export Formats & API Sink

**Branch**: `004-export-formats` | **Date**: 2025-12-04 | **Spec**: [specs/004-export-formats/spec.md](spec.md)
**Input**: Feature specification from `/specs/004-export-formats/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

This feature extends the `rank` command to support JSON and CSV export via a `--format` flag and introduces a generic API sink capability (`--sink-url`) to POST the ranked results to external services (e.g., Trillium Notes). It builds upon the existing `Ranker` service and `RankingRepository`.

## Technical Context

**Language/Version**: Go 1.25.5
**Primary Dependencies**: `net/http` (Stdlib) for Sink POST; `encoding/json` (Stdlib) for JSON export; `encoding/csv` (Stdlib) for CSV export.
**Storage**: SQLite (Read-only access via existing `RankingRepository`).
**Testing**: Go `testing` package + `httptest` for Sink API mocking.
**Target Platform**: CLI (Cross-platform).
**Project Type**: CLI Tool.
**Performance Goals**: Export < 1s for typical datasets. Sink POST latency depends on network.
**Constraints**: JSON/CSV output to stdout; Sink POST is a synchronous blocking call.
**Scale/Scope**: Local processing, outbound HTTP request.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

*   **CLI First**: ✅ Yes, enhances `rank` CLI command.
*   **Modular Architecture**: ✅ Yes, implements `Exporter` interfaces and `Sink` adapter.
*   **Data Integrity**: ✅ Yes, output formats are standard and validated.
*   **TDD**: ✅ Yes, export formatting and sink logic will be tested.
*   **Simplicity**: ✅ Yes, using Stdlib `encoding/json` and `encoding/csv`.

## Project Structure

### Documentation (this feature)

```text
specs/004-export-formats/
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
    └── main.go            # UPDATE: Parse new flags (--format, --sink-url)

internal/
├── core/
│   ├── domain/
│   │   └── interfaces.go  # UPDATE: Add Exporter and Sink interfaces
│   └── service/
│       ├── ranker.go      # UPDATE: Integrate export/sink logic
│       └── ranker_test.go
├── adapter/
│   └── http/              # NEW: Generic HTTP Sink Adapter
│       ├── sink.go
│       └── sink_test.go
├── ui/
│   ├── export.go          # NEW: JSON/CSV Formatting Logic
│   └── export_test.go
```

**Structure Decision**: Standard Go Project Layout. Added `adapter/http` for outbound webhooks and `ui/export.go` for handling different output formats logic.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| None | | |
