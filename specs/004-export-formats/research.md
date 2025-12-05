# Research & Decision Log

## 1. Export Format Handling

**Context**: We need to output ranked repositories in JSON, CSV, or Table formats.
**Question**: Should we use a Strategy Pattern (Interface) or simple switch statement?

**Decision**: **Strategy Pattern (Interface `Exporter`)**

**Rationale**:
*   **Extensibility**: Allows easy addition of future formats (XML, Markdown, YAML) without modifying the core `Ranker` logic.
*   **Testing**: Easier to unit test individual formatters.
*   **Clean Code**: Separates formatting concern from business logic.

**Alternatives Considered**:
*   Switch Statement in `main.go`: Simple but violates Single Responsibility Principle and makes `main` cluttered.

## 2. API Sink Implementation

**Context**: We need to POST data to arbitrary URLs.
**Question**: How to handle headers and payload structure?

**Decision**: **Generic HTTP POST with configurable headers**

**Rationale**:
*   **Flexibility**: Different sinks (Trillium, Slack, Custom) require different headers (Auth, Content-Type).
*   **Standardization**: Sending the standard JSON export format payload ensures consistency. The receiving end is responsible for transformation if needed (or we use a tool like `jq` in between, but built-in sink assumes direct compatibility or generic intake).

## 3. Output vs. Sink Interaction

**Context**: What happens if both `--format` and `--sink-url` are provided?
**Question**: Should they be exclusive?

**Decision**: **Allow both**

**Rationale**:
*   **Utility**: A user might want to see the JSON on stdout AND send it to a webhook simultaneously (e.g., for logging/debugging).
*   **Behavior**: `stdout` gets the formatted output (Table/JSON/CSV). The Sink *always* receives the JSON payload (as it's the most structured/machine-readable standard for APIs). Sending a CSV body to a REST API is uncommon.

**Alternatives Considered**:
*   Exclusive: User can only do one. Restrictive.
*   Sink matches Format: Sending CSV to API if format=csv. Confusing and rare. JSON is the API standard.
