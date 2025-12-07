# Data Model: LLM-Powered Repository Analysis

## Configuration Entities

These entities are persisted in the `config.yaml` file.

### LLMConfig
| Field | Type | Description |
|-------|------|-------------|
| `provider` | string | The identifier of the provider (e.g., "openai", "anthropic", "local"). Currently defaulting to "openai" compatible. |
| `base_url` | string | The API base URL (e.g., "https://api.openai.com/v1", "http://localhost:11434/v1"). |
| `api_key` | string | The authentication key. |
| `model` | string | The model identifier (e.g., "gpt-4o", "llama3"). |
| `max_tokens` | int | (Optional) Max tokens for the response. |

## Ephemeral Entities

These entities are constructed during runtime for the analysis process.

### RepositoryContext
A subset of the `Repository` entity used for LLM prompts to save tokens.

```json
{
  "name": "owner/repo",
  "description": "The description...",
  "language": "Go",
  "stars": 120,
  "forks": 30,
  "tags": ["cli", "tool"]
}
```

### AnalysisRequest
The payload sent to the LLM.

| Field | Type | Description |
|-------|------|-------------|
| `model` | string | The model ID. |
| `messages` | []Message | The chat history/prompt. |

### Message
| Field | Type | Description |
|-------|------|-------------|
| `role` | string | "system" or "user". |
| `content` | string | The text content. |
