# Feature Specification: Configuration Wizard

**Feature Branch**: `005-configuration-wizard`
**Created**: 2025-12-04
**Status**: Draft
**Input**: User description: "Configuration wizard scope Interactive prompts for API URL and Tokens, Persist credentials securely (or at least in a private file like ~/.config/karakeep/config.yaml), Update main.go to load from this config file automatically if flags/env vars are missing."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Interactive Setup (Priority: P1)

As a user, I want to run a setup command that prompts me for my Karakeep URL and API Tokens, so I don't have to manually edit environment variables or configuration files.

**Why this priority**: Reduces friction for new users and simplifies onboarding.

**Independent Test**: Run `karakeep setup` and verify that it prompts for inputs and creates a configuration file.

**Acceptance Scenarios**:

1. **Given** no existing config file, **When** I run `karakeep setup`, **Then** I am prompted for Karakeep URL, Karakeep Token, and GitHub Token.
2. **Given** user inputs valid data, **When** prompts are completed, **Then** a configuration file is created at `~/.config/karakeep/config.yaml` (or OS equivalent) containing the values.
3. **Given** an existing config file, **When** I run `karakeep setup`, **Then** it should prompt for confirmation before overwriting the existing configuration with the new values.

---

### User Story 2 - Automatic Configuration Loading (Priority: P1)

As a user, I want the CLI to automatically load my credentials from the configuration file, so I can run commands like `extract` or `rank` without passing flags every time.

**Why this priority**: This is the core utility of having a persistent configuration.

**Independent Test**: Verify `karakeep extract` runs successfully without flags or env vars, relying solely on the generated config file.

**Acceptance Scenarios**:

1. **Given** a valid `config.yaml`, **When** I run `karakeep extract`, **Then** it uses the credentials from the file.
2. **Given** conflicting inputs (Config File vs Env Var vs Flag), **When** I run a command, **Then** precedence should be: Flag > Env Var > Config File.
3. **Given** a missing or malformed config file, **When** I run a command without flags/env vars, **Then** it should exit with a helpful error suggesting `karakeep setup`.

---

### Edge Cases

- **Permission Denied**: If the tool cannot write to `~/.config/karakeep/`, it should fail gracefully with a permission error.
- **Invalid YAML**: If the config file is corrupted, the tool should warn the user and potentially offer to reset it (or just fail safely).
- **Empty Values**: If user skips optional fields (like GitHub Token) during setup, they should be stored as empty or omitted.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST provide a `setup` CLI command.
- **FR-002**: System MUST interactively prompt the user for `Karakeep URL`, `Karakeep Token`, and `GitHub Token`, masking the input for tokens to prevent visibility.
- **FR-003**: System MUST persist these values to a YAML file located in the user's standard configuration directory (e.g., `~/.config/karakeep/config.yaml` on Linux/macOS), overwriting any existing file after confirmation.
- **FR-004**: System MUST load configuration from this file on startup for all other commands (`extract`, `enrich`, `rank`).
- **FR-005**: System MUST implement configuration precedence: CLI Flags override Environment Variables, which override Config File values.
- **FR-006**: System MUST ensure the configuration directory exists, creating it if necessary.
- **FR-007**: System MUST set appropriate file permissions (0600) on the config file to protect sensitive tokens.

### Key Entities

- **AppConfig**: Struct representing the persistent configuration (URL, Tokens, DB Path).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A new user can configure the tool in under 30 seconds using the wizard.
- **SC-002**: Credentials persist across sessions (verified by running a command in a new terminal window).
- **SC-003**: File permissions on the generated config file are restricted to the user only (0600).

## Clarifications



### Session 2025-12-04



- Q: How should sensitive input (tokens) be handled during the interactive prompt? → A: Masked Input (Hide characters/show asterisks)

- Q: How should the system handle an existing configuration file during setup? → A: Overwrite (Replace the entire file)
