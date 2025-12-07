# Quickstart: LLM Analysis

## 1. Configure the LLM

Before using the analysis feature, you must configure your LLM provider. You can use OpenAI, or any OpenAI-compatible provider (like a local Ollama instance).

```bash
karakeep config llm
```

Follow the prompts:
*   **Base URL**: Enter your provider's URL (e.g., `https://api.openai.com/v1` or `http://localhost:11434/v1`).
*   **API Key**: Enter your key (leave blank if using a local tool that doesn't require one).
*   **Model**: Enter the model name (e.g., `gpt-4o`, `llama3`).

## 2. Run an Analysis

Once configured, you can ask questions about your extracted repositories.

**Example 1: Summarize Go projects**
```bash
karakeep analyze --lang Go "Summarize the top 5 projects and what they do."
```

**Example 2: Find a specific tool**
```bash
karakeep analyze "Which project in my list is a CLI for database management?"
```

## 3. Troubleshooting

*   **"Config not found"**: Run `karakeep config llm` first.
*   **"Context length exceeded"**: Try reducing the number of repos using `--limit` or specific filters.
