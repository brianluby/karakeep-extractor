# Feature Specification: Export Formats & API Sink

**Feature Branch**: `004-export-formats`
**Created**: 2025-12-04
**Status**: Draft
**Input**: User description: "Add export formats (JSON/CSV) and generic API/Webhook sink capability for ranked results"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Standard File Export (Priority: P1)

As a user, I want to export my ranked repository list to JSON or CSV files, so I can open them in spreadsheets or process them with other tools like `jq`.

**Why this priority**: This is the fundamental "utility" requirement for CLI tools—interoperability.

**Independent Test**: Run `karakeep rank --format=json > output.json` and verify valid JSON structure.

**Acceptance Scenarios**:

1. **Given** a populated database, **When** I run `rank --format json`, **Then** the output should be a valid JSON array of repository objects.
2. **Given** a populated database, **When** I run `rank --format csv`, **Then** the output should be valid CSV with headers (Rank, Name, Stars, etc.).
3. **Given** the default behavior (no flag), **When** I run `rank`, **Then** the output should remain the human-readable ASCII table.

---

### User Story 2 - API/Webhook Sink (Priority: P2)

As a user, I want to send the ranked results directly to a custom API endpoint (like Trillium Notes or a webhook), so I can automate my personal knowledge management workflow without manual copy-pasting.

**Why this priority**: Connects the tool to the user's broader ecosystem (Trillium, Zapier, etc.).

**Independent Test**: Run `karakeep rank --sink-url http://localhost:8080` and verify the mock server receives a POST request with the JSON payload.

**Acceptance Scenarios**:

1. **Given** a valid `--sink-url`, **When** I run the command, **Then** the system should POST the ranked results (as JSON) to that URL.
2. **Given** the API returns a success (200 OK), **When** the operation completes, **Then** the CLI should report "Successfully sent to [URL]".
3. **Given** the API fails (500/404), **When** the operation completes, **Then** the CLI should report an error "Failed to send to sink: [Status]".

---

### Edge Cases

- **Large Payloads**: If the list is massive (e.g., 1000 items), the JSON payload might be too large for some webhooks. (Scope limit: send as single batch for now).
- **Invalid URL**: User provides a malformed URL for the sink.
- **Auth**: The sink might require headers (e.g., Authorization). We need a way to pass them.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST support a `--format` flag for the `rank` command with values `table` (default), `json`, `csv`.
- **FR-002**: System MUST support a `--sink-url` flag to define an HTTP POST endpoint.
- **FR-003**: System MUST support a `--sink-header` flag (repeatable) to pass custom headers (e.g., `Authorization: Bearer ...`) to the sink.
- **FR-004**: When `--sink-url` is provided, the system MUST POST the result set as a JSON body to the target URL.
- **FR-005**: JSON output MUST include all available fields: `rank`, `repo_id`, `url`, `stars`, `forks`, `last_pushed_at`, `description`, `language`.
- **FR-006**: CSV output MUST include a header row and comma-separated values for the same fields.

### Key Entities

- **ExportPayload**: The structure sent to the sink or output as JSON. List of `RankedRepo`.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: JSON output is parseable by standard tools (`jq`) without errors.
- **SC-002**: CSV output opens correctly in Excel/Numbers with columns aligned.
- **SC-003**: Integration with a local mock server confirms correct HTTP method (POST), headers, and payload body.