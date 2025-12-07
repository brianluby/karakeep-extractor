# Research: LLM-Powered Repository Analysis

**Status**: Complete
**Date**: 2025-12-06

## Decisions

### 1. LLM Client Implementation
**Decision**: Implement a lightweight, custom HTTP client using Go's standard `net/http` library rather than importing a third-party SDK.
**Rationale**: 
- **Constitution Compliance**: Adheres to the "Minimize external dependencies" principle.
- **Simplicity**: We only need a specific subset of functionality (Chat Completions API), which is easy to implement with a few structs and an HTTP request.
- **Control**: Allows us to easily handle different providers (OpenAI, Anthropic, generic local endpoints) by just changing the request body/headers if needed, without fighting an SDK's abstraction.

### 2. Configuration Storage
**Decision**: Extend the existing `Config` struct in `internal/config/config.go` and the YAML file structure.
**Rationale**:
- **Consistency**: Keeps all application configuration in one place (`~/.config/karakeep/config.yaml`).
- **Ease of Use**: Users already use this file; adding a new section `llm:` is intuitive.

### 3. Prompt Engineering & Context Management
**Decision**: Use JSON format to serialize repository data within the prompt.
**Rationale**:
- **Reliability**: LLMs are generally very good at parsing and understanding JSON.
- **Structure**: Easy to select specific fields (Name, Description, Language, Stars) to minimize token usage compared to raw text dumps.
- **Limit**: Implement a hard limit on the number of repositories (e.g., top 50 by ranking) or a token estimation to prevent overflowing the context window.

## Alternatives Considered

### Use `sashabaranov/go-openai`
- **Pros**: Complete API coverage, community supported.
- **Cons**: Adds a dependency, primarily focused on OpenAI (though compatible with others). Overkill for just one endpoint.
- **Result**: Rejected in favor of stdlib.

### Local LLM via Bindings (e.g., `go-llama.cpp`)
- **Pros**: Runs offline, fast.
- **Cons**: Complex build process (CGO), huge binary size, hardware dependency.
- **Result**: Rejected. We will support local LLMs via their HTTP API servers (e.g., Ollama, LM Studio) which use the standard OpenAI-compatible endpoints.

## Open Questions Resolved

- **Provider Support**: We will initially support OpenAI-compatible APIs (which covers OpenAI, Groq, DeepSeek, and local tools like Ollama/LM Studio). This provides the widest compatibility with the least code.
