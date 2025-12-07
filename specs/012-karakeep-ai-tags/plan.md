# Implementation Plan: Karakeep AI Tags Extraction

**Branch**: `012-karakeep-ai-tags` | **Date**: 2025-12-06 | **Spec**: [specs/012-karakeep-ai-tags/spec.md](spec.md)
**Input**: Feature specification from `/specs/012-karakeep-ai-tags/spec.md`

## Summary

Extract AI-generated tags associated with bookmarks from the Karakeep API and store them in a normalized SQLite schema. This involves updating the `RawBookmark` and `ExtractedRepo` entities, implementing database migrations for `tags` and `repo_tags` tables, modifying the extraction service to parse and persist tags, and updating the ranking repository to support filtering by these tags.

## Technical Context

**Language/Version**: Go 1.25.5
**Primary Dependencies**: `github.com/mattn/go-sqlite3` (Existing).
**Storage**: SQLite.
**Testing**: Unit tests for tag parsing and database persistence.
**Target Platform**: CLI.
**Project Type**: Single CLI binary.
**Performance Goals**: Tag filtering < 2s.
**Constraints**: Backward compatibility for existing databases (migrations).

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- [x] **CLI First**: Enhances `extract` and `rank` CLI commands.
- [x] **Modular Architecture**: Changes restricted to Data/Repo layer and Extractor service.
- [x] **Data Integrity**: Normalizes tags to prevent duplication.
- [x] **TDD**: Tests will cover new DB logic.
- [x] **Simplicity**: Uses standard SQL relationships.

## Project Structure

### Documentation (this feature)

```text
specs/012-karakeep-ai-tags/
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
│   │   ├── domain/      # Update RawBookmark, ExtractedRepo
│   ├── adapter/
│   │   ├── sqlite/      # Update Migrations, Repository (Save/GetRankedRepos)
│   └── core/
│       └── service/     # Update Extractor to handle tags
```

**Structure Decision**: Option 1: Single project (DEFAULT).

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| N/A       |            |                                     |