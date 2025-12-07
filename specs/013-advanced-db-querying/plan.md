# Implementation Plan: Advanced Local Database Querying

**Branch**: `013-advanced-db-querying` | **Date**: 2025-12-06 | **Spec**: [specs/013-advanced-db-querying/spec.md](spec.md)
**Input**: Feature specification from `/specs/013-advanced-db-querying/spec.md`

## Summary

Add a powerful `query` CLI command to filter the local SQLite database using criteria like star counts, fork counts, and extraction dates. This involves creating a new `QueryFilter` domain entity, adding a `QueryRepos` method to the `RepoRepository` interface (and SQLite implementation), and parsing CLI flags to construct dynamic SQL queries.

## Technical Context

**Language/Version**: Go 1.25.5
**Primary Dependencies**: `github.com/mattn/go-sqlite3` (Existing).
**Storage**: SQLite.
**Testing**: Unit tests for query building and flag parsing.
**Target Platform**: CLI.
**Project Type**: Single CLI binary.
**Performance Goals**: Queries < 500ms.
**Constraints**: Dynamic SQL building must use parameters to prevent injection.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- [x] **CLI First**: Adds new `query` command.
- [x] **Modular Architecture**: Logic resides in Repo adapter and Service layer.
- [x] **Data Integrity**: Safe SQL usage.
- [x] **TDD**: Tests for query builder.
- [x] **Simplicity**: Reuses existing infrastructure.

## Project Structure

### Documentation (this feature)

```text
specs/013-advanced-db-querying/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output
└── tasks.md             # Phase 2 output
```

### Source Code (repository root)

```text
src/
├── internal/
│   ├── core/
│   │   ├── domain/      # Add QueryFilter struct
│   ├── adapter/
│   │   ├── sqlite/      # Implement QueryRepos
│   └── ui/
│       └── filter.go    # Helper to parse CLI flags into QueryFilter
```

**Structure Decision**: Option 1: Single project (DEFAULT).

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| N/A       |            |                                     |