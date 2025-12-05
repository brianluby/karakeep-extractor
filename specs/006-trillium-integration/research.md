# Research & Decision Log

## 1. Trillium Content Format

**Context**: Trillium notes store content as HTML (rich text) or plain text.
**Question**: Should we send raw HTML or Markdown?

**Decision**: **Markdown (converted to HTML if necessary, or plain text)**

**Rationale**:
*   **Simplicity**: Generating a Markdown table string is easier than constructing valid HTML DOM structure in Go.
*   **Trillium Support**: Trillium supports Markdown-like syntax or code blocks. However, standard "Text" notes in Trillium are HTML.
*   **Refinement**: We will create a "Code" note type (language: markdown) OR a "Text" note where we inject a basic HTML table.
*   *Correction*: Trillium API `create-note` allows `type: 'text'`. If we want a nicely formatted table, HTML is safer for the 'text' type. But a "Markdown" code note is cleaner for raw data.
*   **Final Decision**: Create a **Text** note and render a simple **HTML Table**. It looks better in Trillium's WYSIWYG editor than raw markdown text.

**Alternatives Considered**:
*   Markdown Note: Good for copy-pasting out, but less readable in-app.
*   HTML Table: Best for in-app viewing.

## 2. Trillium Client Implementation

**Context**: Need to talk to ETAPI.
**Question**: Use a library or custom client?

**Decision**: **Custom Client (using `net/http`)**

**Rationale**:
*   **Minimal dependencies**: The API is simple (`POST /etapi/create-note`). No need for a full SDK.
*   **Control**: We only need one specific endpoint.

## 3. Sink Integration Strategy

**Context**: How to fit into `Ranker`.
**Question**: Generic Sink vs Dedicated Flag?

**Decision**: **Dedicated Flag `--sink-trillium` that instantiates a `TrilliumSink`**

**Rationale**:
*   **UX**: Easier than asking user to type `--sink-url ... --sink-header ...` for a known integration.
*   **Config**: Can pull credentials from `config.yaml` specifically for Trillium.
*   **Polymorphism**: `TrilliumSink` will implement the `domain.Sink` interface, so `Ranker` doesn't need to change (it just receives a `Sink`).

## 4. Setup Flow UX

**Context**: User has multiple potential integrations.
**Question**: Ask for all or opt-in?

**Decision**: **Opt-in Prompt**

**Rationale**:
*   **User Experience**: "Do you want to configure Trillium? (y/N)" prevents asking for tokens the user doesn't have.
*   **Future-proofing**: As we add Obsidian/Notion, this pattern scales better than a monolithic list of questions.
