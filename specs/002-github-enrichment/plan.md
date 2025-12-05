# Implementation Plan: GitHub Repository Enrichment

**Branch**: `002-github-enrichment` | **Date**: 2025-12-04 | **Spec**: [specs/002-github-enrichment/spec.md](spec.md)
**Input**: Feature specification from `/specs/002-github-enrichment/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

This feature implements the "Enrichment" phase of the Karakeep pipeline. It scans the existing `extracted_repos` SQLite table for GitHub URLs, fetches metadata (Stars, Forks, Description, PushedAt, Language) from the GitHub REST API using a documented PAT or unauthenticated access, and updates the database rows. It prioritizes rate-limit safety (fail-fast) and data persistence.

## Technical Context

**Language/Version**: Go 1.25.5
**Primary Dependencies**: `net/http` (Standard Library) for API interaction; `encoding/json` for parsing; `github.com/mattn/go-sqlite3` (Existing) for storage.
**Storage**: SQLite (Extension of existing `extracted_repos` table).
**Testing**: Go `testing` package + `net/http/httptest` for mocking GitHub API.
**Target Platform**: CLI (Cross-platform: macOS, Linux, Windows).
**Project Type**: CLI Tool.
**Performance Goals**: Process ~50 repos < 10s. Concurrent fetching (worker pool pattern).
**Constraints**: GitHub API Rate Limits (60/hr unauth, 5000/hr auth). strict "Exit/Fail Fast" behavior on limit breached.
**Scale/Scope**: Single user, local database. Batch processing.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

*   **CLI First**: ✅ Yes, adds `enrich` command or flag to existing CLI.
*   **Modular Architecture**: ✅ Yes, adds `internal/adapter/github` and updates `internal/core/service`. Decoupled from extraction.
*   **Data Integrity**: ✅ Yes, handles 404s and rate limits explicitly. Updates DB transactionally or row-by-row.
*   **TDD**: ✅ Yes, `httptest` will be used to simulate GitHub responses before implementation.
*   **Simplicity**: ✅ Yes, using `net/http` instead of heavy client libraries. Inline SQL migrations.
*   **Tech Constraints**: ✅ Yes, Go 1.25.5, minimal deps.

## Project Structure

### Documentation (this feature)

```text
specs/002-github-enrichment/
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
    └── main.go            # Wiring up the new Enricher service

internal/
├── adapter/
│   ├── github/            # NEW: GitHub API Client implementation
│   │   ├── client.go
│   │   └── client_test.go
│   └── sqlite/
│       ├── repository.go  # UPDATE: Add UpdateRepoStats method
│       └── repository_test.go
├── core/
│   ├── domain/
│   │   └── entities.go    # UPDATE: Add RepoStats fields to Entity
│   ├── port/              # NEW/UPDATE: Define interfaces (Enricher, Repository additions)
│   │   └── ports.go       # (Assuming ports are here or in domain/interfaces.go)
│   └── service/
│       ├── enricher.go    # NEW: Orchestration logic (Rate limits, Worker pool)
│       └── enricher_test.go
└── config/
    └── config.go          # UPDATE: Add GitHubToken field
```

**Structure Decision**: Standard Go Project Layout (cmd/internal) as established in `001-karakeep-extract`.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| None | | |
