# Implementation Plan: LLM-Powered Repository Analysis

**Branch**: `010-llm-repo-analysis` | **Date**: 2025-12-06 | **Spec**: [specs/010-llm-repo-analysis/spec.md](spec.md)
**Input**: Feature specification from `/specs/010-llm-repo-analysis/spec.md`

## Summary

Integrate with an LLM (Large Language Model) API to perform advanced analysis on the extracted repository data. This feature enables users to query their local database of repositories using natural language prompts (e.g., "Summarize these projects"). It involves adding a new `llm` configuration section, an `analyze` CLI command, and a lightweight HTTP client to interact with OpenAI-compatible APIs (including local models).

## Technical Context

**Language/Version**: Go 1.25.5
**Primary Dependencies**: `net/http` (Standard Lib), `encoding/json` (Standard Lib). No external SDKs.
**Storage**: SQLite (for repository data), `config.yaml` (for LLM credentials).
**Testing**: Go standard testing package.
**Target Platform**: CLI (Cross-platform).
**Project Type**: Single CLI binary.
**Performance Goals**: Interactive response times (<15s for API calls).
**Constraints**: Must respect LLM context window limits; prioritize standard library.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- [x] **CLI First**: Adds `karakeep config llm` and `karakeep analyze` commands.
- [x] **Modular Architecture**: LLM client will be a separate package/service, decoupled from the core extraction logic.
- [x] **Data Integrity**: Will validate API responses and handle errors gracefully.
- [x] **TDD**: Client and Prompt generation will be unit tested.
- [x] **Simplicity**: Using stdlib `net/http` instead of heavy SDKs.

## Project Structure

### Documentation (this feature)

```text
specs/010-llm-repo-analysis/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output
│   └── cli-commands.md
└── tasks.md             # Phase 2 output
```

### Source Code (repository root)

```text
src/
├── internal/
│   ├── config/          # Update Config struct
│   ├── core/
│   │   ├── domain/      # Update entities (if needed)
│   │   └── service/
│   │       └── analysis/ # New service for Analysis logic
│   ├── adapter/
│   │   └── llm/         # New adapter for LLM API interaction
│   └── ui/
│       └── prompt.go    # UI prompts for config wizard
```

**Structure Decision**: Option 1: Single project (DEFAULT). We are adding a new adapter (`llm`) and a new service (`analysis`) within the existing `internal` structure.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| N/A       |            |                                     |