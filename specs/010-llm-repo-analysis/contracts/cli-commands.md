# CLI Contract: LLM Analysis

## Commands

### `karakeep config llm`

Interactive wizard to configure LLM settings.

**Usage:**
```bash
karakeep config llm
```

**Interaction Flow:**
1. System checks for existing config.
2. Prompts for `Provider Base URL` (default: https://api.openai.com/v1).
3. Prompts for `API Key` (masked input).
4. Prompts for `Model Name` (default: gpt-4o).
5. Validates connection (optional "ping").
6. Saves to config file.

### `karakeep analyze`

Analyze repositories using the configured LLM.

**Usage:**
```bash
karakeep analyze [flags] <query>
```

**Arguments:**
*   `<query>`: The natural language question/instruction (e.g., "Summarize these projects").

**Flags:**
*   `--lang <string>`: Filter repositories by language (e.g., "Go").
*   `--limit <int>`: Limit the number of repositories to include (default: 50).
*   `--tag <string>`: Filter by tag (if available).

**Example:**
```bash
karakeep analyze --lang Go "Which of these libraries is best for HTTP requests?"
```

**Output:**
*   Streamed or final text response from the LLM.
*   Error message if config is missing or API fails.
