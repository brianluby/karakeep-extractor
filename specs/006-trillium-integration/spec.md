# Feature Specification: Trillium Notes Integration

**Feature Branch**: `006-trillium-integration`
**Created**: 2025-12-04
**Status**: Draft
**Input**: User description: "Let's add Trillium Integration capability. Trillium Notes: Direct API integration to store the ranked list in markdown format as a page in trillium notes. here is the api connectivity details https://github.com/TriliumNext/Trilium/blob/main/docs/Developer%20Guide/Developer%20Guide/Architecture/APIs.md specifically the ETAPI (external api)"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Export to Trillium (Priority: P1)

As a user, I want to send my ranked repository list directly to my Trillium Notes instance as a new note, so I can archive my research without manual copy-pasting.

**Why this priority**: Seamless integration with PKM tools is a key value prop for "knowledge collection".

**Independent Test**: Run `karakeep rank --sink-trillium` (with config) and verify a new note appears in Trillium with the correct content.

**Acceptance Scenarios**:

1. **Given** valid Trillium credentials in config, **When** I run `karakeep rank --sink-trillium`, **Then** a new note is created in Trillium containing the ranked list in Markdown format.
2. **Given** the list is updated, **When** I run the command again, **Then** a *new* note is created (default behavior to preserve history).
3. **Given** invalid credentials, **When** I run the command, **Then** the CLI reports an authentication error from the Trillium API.

---

### User Story 2 - Configure Integrations in Setup (Priority: P2)

As a user, I want the setup wizard to offer me a choice of integrations to configure (like Trillium), so I only have to provide details for the services I actually use.

**Why this priority**: Prevents setup fatigue by making optional integrations opt-in. Sets the stage for adding Obsidian/Notion later.

**Independent Test**: Run `karakeep setup` and verify it asks "Do you want to configure Trillium Notes? [y/N]" before prompting for credentials.

**Acceptance Scenarios**:

1. **Given** I run `karakeep setup`, **When** I finish the core config, **Then** the wizard asks "Configure Trillium Integration? [y/N]".
2. **Given** I answer "Yes", **Then** I am prompted for "Trillium Instance URL" and "Trillium ETAPI Token".
3. **Given** I answer "No", **Then** the wizard skips Trillium configuration and finishes.

---

### Edge Cases

- **API Version Mismatch**: Trillium API changes. (We target stable ETAPI `v1`).
- **Network unreachable**: Trillium instance is self-hosted and offline. Report "Connection Refused".
- **Parent Note**: Where should the notes go? Root? A specific "Inbox"? (Default: Root or "Inbox" if found, else Root).

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST support a `--sink-trillium` flag for the `rank` command to trigger this specific sink.
- **FR-002**: System MUST implement a Trillium Client that communicates via ETAPI (`/etapi/create-note`).
- **FR-003**: System MUST format the ranked repository list as a Markdown table (or list) suitable for Trillium's content format.
- **FR-004**: System MUST persist Trillium configuration (`TrilliumURL`, `TrilliumToken`) in `config.yaml` via the `setup` command.
- **FR-005**: System MUST create a note of type `text` with the title "GitHub Rankings - [Timestamp]".
- **FR-006**: System MUST handle HTTP errors from Trillium gracefully.
- **FR-007**: The `setup` command MUST introduce an "Integrations" step that asks the user for opt-in before prompting for integration-specific credentials.

### Key Entities

- **TrilliumConfig**: Struct for URL/Token.
- **TrilliumNote**: Payload struct for creating a note (`title`, `type`, `content`).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A user can push a ranked list to Trillium in one command after setup.
- **SC-002**: The created note renders correctly in Trillium (valid Markdown).
- **SC-003**: Setup wizard successfully saves and loads Trillium credentials.