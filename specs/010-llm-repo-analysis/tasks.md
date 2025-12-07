# Tasks: LLM-Powered Repository Analysis

**Feature**: LLM-Powered Repository Analysis
**Spec**: [specs/010-llm-repo-analysis/spec.md](spec.md)
**Plan**: [specs/010-llm-repo-analysis/plan.md](plan.md)

## Phase 1: Setup
*Goal: Initialize project structure for the new feature.*

- [x] T001 Create directory structure for LLM adapter and service `internal/adapter/llm`, `internal/core/service/analysis`
- [x] T002 Create `internal/core/domain/llm.go` to define `LLMConfig`, `AnalysisRequest`, `Message`, and `RepositoryContext` structs

## Phase 2: Foundational
*Goal: Update configuration to support LLM settings. Blocking US1 & US2.*

- [x] T003 Update `Config` struct in `internal/config/config.go` to include `LLMConfig` fields (Provider, BaseURL, APIKey, Model)
- [x] T004 Update `internal/config/loader.go` to parse the new `llm` YAML section
- [x] T005 [TEST] Add unit tests for saving/loading LLM configuration in `internal/config/config_test.go`

## Phase 3: User Story 1 - Configure LLM Connection
*Goal: Allow users to set up their LLM provider credentials via CLI.*

**Independent Test**: Run `karakeep config llm`, enter credentials, and verify they are saved to `~/.config/karakeep/config.yaml`.

- [x] T006 [US1] Add `config` subcommand and `llm` sub-subcommand logic to `main` function in `cmd/extractor/main.go`
- [x] T007 [US1] Implement `runConfigLLM` function in `cmd/extractor/main.go` (or separate file) using `internal/ui/prompt` to ask for LLM settings
- [x] T008 [US1] Implement logic to save the updated config using `config.SaveConfig` in `cmd/extractor/main.go`

## Phase 4: User Story 2 - Analyze Repository with Natural Language
*Goal: Enable users to query repository data using natural language.*

**Independent Test**: Run `karakeep analyze --limit 5 "Summarize these repos"` and verify output contains a text summary.

- [x] T009 [US2] [TEST] Create `internal/adapter/llm/client_test.go` to test `SendMessage` with a mock HTTP server
- [x] T010 [US2] Implement `Client` struct and `SendMessage` method in `internal/adapter/llm/client.go` using `net/http`
- [x] T011 [US2] Implement prompt construction logic (JSON serialization of `RepositoryContext`) in `internal/core/service/analysis/prompt.go`
- [x] T012 [US2] [TEST] Add unit tests for prompt construction in `internal/core/service/analysis/prompt_test.go`
- [x] T013 [US2] Implement `Analyze` method in `internal/core/service/analysis/service.go` to coordinate fetching repos, building prompts, and calling LLM
- [x] T014 [US2] Add `analyze` subcommand and flags (`--lang`, `--limit`, `--tag`) to `cmd/extractor/main.go`
- [x] T015 [US2] Implement `runAnalyze` in `cmd/extractor/main.go` to wire up the service and print results

## Phase 5: Polish & Cross-Cutting
*Goal: Ensure robust error handling and documentation.*

- [x] T016 Add user-friendly error messages for network timeouts and invalid API keys in `internal/adapter/llm/client.go` and `cmd/extractor/main.go`
- [x] T017 Update `README.md` with documentation for `karakeep config llm` and `karakeep analyze` usage

## Implementation Strategy
- **MVP**: Complete Phase 1, 2, and 3 to get configuration working. Then Phase 4 to get basic analysis (even if just 1 repo) working.
- **Incremental**: First make `analyze` work with a hardcoded prompt, then add the dynamic context builder.

## Dependencies
- **Phase 2** blocks **Phase 3** and **Phase 4** (Config structs needed).
- **Phase 3** blocks **Phase 4** runtime (need config to run analysis), but code can be written in parallel.
- **T010 (Client)** and **T011 (Prompt)** can be done in parallel.

## Parallel Execution Examples
- **US2**: Developer A implements `internal/adapter/llm` (T010) while Developer B implements `internal/core/service/analysis/prompt.go` (T011).