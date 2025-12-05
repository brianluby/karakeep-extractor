<!-- Sync Impact Report
Version: 1.0.0 -> 1.1.0
Modified Principles: Technical Constraints (Language updated to Go)
Added Sections: None
Removed Sections: None
Templates requiring updates: 
- .specify/templates/tasks-template.md (⚠ pending: "Tests OPTIONAL" contradicts Constitution TDD principle)
Follow-up TODOs: None
-->
# Karakeep Extractor Constitution

## Core Principles

### I. CLI First
All core functionality must be exposed via the Command Line Interface (CLI). The tool operates primarily as a CLI utility, adhering to standard conventions (exit codes, stdout/stderr separation, help flags). Visual interfaces (like Web UI) are secondary and must build upon the CLI/Core logic, not replace it.

### II. Modular Architecture
The system is divided into distinct, loosely coupled phases: Extraction, Enrichment, and Ranking. Components (Source, Processor, Enricher, Output) must be independent to allow for isolated testing and future replacement.

### III. Data Integrity & Validation
Data extracted from external sources (Karakeep, GitHub) is inherently unstable. The system must rigorously validate all inputs. Missing or malformed data should be handled gracefully—logging warnings rather than crashing the entire pipeline whenever possible.

### IV. Test-Driven Development (TDD)
Critical logic, especially parsers, enrichers, and ranking algorithms, must be covered by tests. Tests should be written to define the expected behavior before the implementation details are finalized.

### V. Simplicity & Focus
Prioritize the core "Extraction -> Enrichment -> Ranking" loop. Avoid premature optimization or feature bloat (e.g., complex UIs) until the core value proposition is solid and reliable (v1).

## Technical Constraints

*   **Language:** Go (Golang) 1.25+ (Preferred for concurrency and performance).
*   **Dependencies:** Minimize external dependencies to reduce the surface area for security issues and maintenance.
*   **Output:** Must support structured output (JSON/CSV) to enable piping to other tools.

## Development Workflow

*   **Branching:** Feature branches merged via Pull Request.
*   **Versioning:** Semantic Versioning (MAJOR.MINOR.PATCH).
    *   MAJOR: Breaking CLI changes or data format changes.
    *   MINOR: New features (e.g., new source or enrichment provider).
    *   PATCH: Bug fixes.

## Governance

This constitution supersedes other informal practices.
*   **Amendments:** Require a documented update to this file and a version bump.
*   **Compliance:** All PRs must be checked against these principles.
*   **Guidance:** Use `.specify/templates/plan-template.md` for feature planning to ensure alignment.

**Version**: 1.1.0 | **Ratified**: 2025-12-04 | **Last Amended**: 2025-12-04
