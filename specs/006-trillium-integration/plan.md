# Implementation Plan: Trillium Notes Integration

**Branch**: `006-trillium-integration` | **Date**: 2025-12-04 | **Spec**: [specs/006-trillium-integration/spec.md](spec.md)
**Input**: Feature specification from `/specs/006-trillium-integration/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

This feature implements a specialized API sink for Trillium Notes. It adds a `--sink-trillium` flag to the `rank` command, which formats the repository list as a Markdown table and creates a new note in Trillium via its ETAPI. It also updates the configuration wizard to support optional integration setup for Trillium credentials.

## Technical Context

**Language/Version**: Go 1.25.5
**Primary Dependencies**: `net/http` (Stdlib) for API calls.
**Storage**: Configuration update (YAML) to store Trillium credentials.
**Testing**: Go `testing` package + `httptest` for mocking Trillium ETAPI.
**Target Platform**: CLI (Cross-platform).
**Project Type**: CLI Tool.
**Performance Goals**: API call < 2s (network dependent).
**Constraints**: Must handle Trillium's specific JSON payload format for note creation.
**Scale/Scope**: Single user integration.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

*   **CLI First**: ✅ Yes, adds CLI flag and setup prompts.
*   **Modular Architecture**: ✅ Yes, implements a new `Sink` strategy (`TrilliumSink`) separate from core logic.
*   **Data Integrity**: ✅ Yes, validates API response codes.
*   **TDD**: ✅ Yes, client logic will be tested with mocks.
*   **Simplicity**: ✅ Yes, reusing existing HTTP client patterns.

## Project Structure

### Documentation (this feature)

```text
specs/006-trillium-integration/
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
    └── main.go            # UPDATE: Wire up --sink-trillium flag and setup prompts

internal/
├── adapter/
│   └── trillium/          # NEW: Trillium API Client / Sink
│       ├── client.go
│       ├── client_test.go
│       └── sink.go        # Implements domain.Sink
├── config/
│   └── config.go          # UPDATE: Add Trillium fields
├── ui/
│   ├── prompt.go          # (Existing)
│   └── markdown.go        # NEW: Markdown table formatter (if not reusing TableRenderer)
```

**Structure Decision**: Adding `internal/adapter/trillium` to encapsulate the specific API logic, keeping it separate from the generic `http` sink.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| None | | |
