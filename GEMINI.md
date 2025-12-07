# karakeep-extractor Development Guidelines

Auto-generated from all feature plans. Last updated: 2025-12-04

## Active Technologies
- Go 1.25.5 + `net/http` (Standard Library) for API interaction; `encoding/json` for parsing; `github.com/mattn/go-sqlite3` (Existing) for storage. (002-github-enrichment)
- SQLite (Extension of existing `extracted_repos` table). (002-github-enrichment)
- Go 1.25.5 + `net/http` (Stdlib), `encoding/json` (Stdlib), `github.com/mattn/go-sqlite3` (Existing) (003-github-ranking)
- SQLite (Read-only access for this feature) (003-github-ranking)
- Go 1.25.5 + `net/http` (Stdlib) for Sink POST; `encoding/json` (Stdlib) for JSON export; `encoding/csv` (Stdlib) for CSV export. (004-export-formats)
- SQLite (Read-only access via existing `RankingRepository`). (004-export-formats)
- Go 1.25.5 + `gopkg.in/yaml.v3` (for YAML parsing - standard in Go ecosystem for config), `bufio` (Stdlib) for interactive prompts. (005-configuration-wizard)
- YAML file (`~/.config/karakeep/config.yaml`). (005-configuration-wizard)
- Go 1.25.5 + `net/http` (Stdlib) for API calls. (006-trillium-integration)
- Configuration update (YAML) to store Trillium credentials. (006-trillium-integration)
- Go 1.25.5 + `github.com/mattn/go-sqlite3` (Existing). (007-filter-by-tag)
- SQLite (Read-only query modification). (007-filter-by-tag)
- Markdown (for docs), Go 1.25.5 (for CLI help text updates). + None. (008-documentation-sprint)
- N/A (Purely UI feature) (009-optional-tui-status)
- Go 1.25.5 + `net/http` (Standard Lib), `encoding/json` (Standard Lib). No external SDKs. (010-llm-repo-analysis)
- SQLite (for repository data), `config.yaml` (for LLM credentials). (010-llm-repo-analysis)
- Go 1.25.5 + `github.com/charmbracelet/bubbletea` (Existing TUI framework). (011-real-time-tui-stats)
- N/A (Transient UI state). (011-real-time-tui-stats)

- Go 1.25.5 + `net/http` (Stdlib), `github.com/mattn/go-sqlite3` (SQLite Driver) (001-karakeep-extract)

## Project Structure

```text
src/
tests/
```

## Commands

# Add commands for Go 1.25.5

## Code Style

Go 1.25.5: Follow standard conventions

## Recent Changes
- 013-advanced-db-querying: Added Go 1.25.5 + `github.com/mattn/go-sqlite3` (Existing).
- 013-advanced-db-querying: Added [if applicable, e.g., PostgreSQL, CoreData, files or N/A]
- 012-karakeep-ai-tags: Added Go 1.25.5 + `github.com/mattn/go-sqlite3` (Existing).


<!-- MANUAL ADDITIONS START -->
<!-- MANUAL ADDITIONS END -->

<!-- BACKLOG.MD MCP GUIDELINES START -->

<CRITICAL_INSTRUCTION>

## BACKLOG WORKFLOW INSTRUCTIONS

This project uses Backlog.md MCP for all task and project management activities.

**CRITICAL GUIDANCE**

- If your client supports MCP resources, read `backlog://workflow/overview` to understand when and how to use Backlog for this project.
- If your client only supports tools or the above request fails, call `backlog.get_workflow_overview()` tool to load the tool-oriented overview (it lists the matching guide tools).

- **First time working here?** Read the overview resource IMMEDIATELY to learn the workflow
- **Already familiar?** You should have the overview cached ("## Backlog.md Overview (MCP)")
- **When to read it**: BEFORE creating tasks, or when you're unsure whether to track work

These guides cover:
- Decision framework for when to create tasks
- Search-first workflow to avoid duplicates
- Links to detailed guides for task creation, execution, and completion
- MCP tools reference

You MUST read the overview resource to understand the complete workflow. The information is NOT summarized here.

</CRITICAL_INSTRUCTION>

<!-- BACKLOG.MD MCP GUIDELINES END -->
