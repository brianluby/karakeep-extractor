# Feature Specification: Optional TUI Status

**Feature Branch**: `009-optional-tui-status`  
**Created**: 2025-12-05  
**Status**: Draft  
**Input**: User description: "We should build and optional tui that has a some sort of indicator of what it is doing."

## Clarifications

### Session 2025-12-05
- Q: Which TUI library should be used? → A: **Bubble Tea** (`charmbracelet/bubbletea`) for robust state management and modern UI.
- Q: How to decouple TUI from core logic? → A: **Observer Pattern** (Reporter Interface). Core services will accept a `ProgressReporter` interface to send events, allowing swappable UI (TUI vs. Stdout).

- Q: How should the main program loop be structured given a TUI library? → A: **UI-Driven**: The `tea.Program` will own the main loop, with work performed in goroutines communicating via `tea.Msg`s.
- Q: What is the behavior for the error/log display area within the TUI? → A: **Static Tail**: Display the last N (e.g., 3-5) lines of errors/status messages in a dedicated area.
- Q: For Extraction TUI, is a determinate progress bar (requiring a known total upfront) necessary, or is an indeterminate spinner/counter sufficient? → A: **Indeterminate (Spinner/Counter)**: Karakeep API typically returns paginated results without an upfront total count.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - TUI for Enrichment Process (Priority: P1)

As a user running the enrichment command, I want to see a visual interface with a progress bar so that I can estimate how much time remains and monitor the status of individual repository processing.

**Why this priority**: The enrichment process is the longest-running operation (fetching data for potentially hundreds of repos). A visual indicator provides the most value here by reducing user uncertainty.

**Independent Test**: Run `karakeep-extractor enrich --tui` with a set of bookmarks. Verify that a progress bar appears, updates as repos are processed, and the application exits cleanly upon completion.

**Acceptance Scenarios**:

1. **Given** the user has bookmarks to enrich, **When** they run `karakeep-extractor enrich --tui`, **Then** the terminal screen clears and displays a progress bar indicating the percentage of completed repositories.
2. **Given** the TUI is running, **When** a repository is successfully enriched, **Then** the progress bar increments and a status message updates (e.g., "Enriched: owner/repo").
3. **Given** the TUI is running, **When** a rate limit or error occurs, **Then** the error is displayed in a dedicated log area within the TUI without disrupting the progress bar.
4. **Given** the user runs `karakeep-extractor enrich` (without flag), **Then** the standard text log output is displayed, and no TUI is initialized.

---

### User Story 2 - TUI for Extraction Process (Priority: P2)

As a user running the extraction command, I want to see a visual status indicator so I know the tool is working and how many bookmarks have been processed.

**Why this priority**: Extraction can also take time depending on the number of bookmark pages. Visual feedback confirms the connection is active.

**Independent Test**: Run `karakeep-extractor extract --tui`. Verify that a spinner or counter updates as pages are fetched.

**Acceptance Scenarios**:

1. **Given** valid Karakeep credentials, **When** the user runs `karakeep-extractor extract --tui`, **Then** a TUI appears showing the current page being fetched and the count of GitHub links found.
2. **Given** the extraction finishes, **When** the process is complete, **Then** the TUI displays a summary (e.g., "Extracted 150 repos") and prompts to exit or exits automatically.

---

### Edge Cases

- **Terminal Resize**: The TUI should handle terminal window resizing without crashing or garbling the output.
- **Zero Data**: If there are no repos to enrich or extract, the TUI should display a "Nothing to do" message and exit gracefully.
- **SIGINT (Ctrl+C)**: If the user interrupts the process, the TUI should clean up the terminal (restore cursor, clear special formatting) before exiting.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST support a `--tui` flag for `extract` and `enrich` commands.
- **FR-002**: When TUI mode is enabled, the system MUST hide standard standard log output (stdout/stderr) and replace it with the visual interface.
- **FR-003**: The Enrichment TUI MUST display a progress bar representing `processed / total` repositories.
- **FR-004**: The TUI MUST display the name of the repository currently being processed.
- **FR-005**: The TUI MUST provide a distinct visual indication or log area for errors (e.g., "404 Not Found", "Rate Limit Exceeded").
- **FR-006**: The system MUST restore the terminal to its original state (cursor visible, no artifacts) upon exit or failure.
- **FR-007**: If the `--tui` flag is NOT provided, the system MUST behave exactly as it did previously (standard text logs).
- **FR-008**: The TUI implementation MUST use `charmbracelet/bubbletea` for the render loop and event handling.

### Key Entities

- **TUI State**: Holds the current progress (count, total), current item description, and a list of recent logs/errors.
- **Progress Model**: Represents the percentage completion of the current batch operation.
- **ProgressReporter Interface**: Abstraction for reporting status events (Start, Progress, Error, Finish) to the active UI (Text or TUI).

### Assumptions

- The user is running the application in a standard terminal environment that supports ANSI escape codes (for TUI rendering).
- The TUI mode is optional; the application must remain fully functional in "headless" environments (CI/CD) without this flag.
- No specific accessibility compliance (e.g., screen readers) is required for this MVP, though standard text logs serve as the accessible fallback.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can determine the exact progress (e.g., "50%") of an enrichment task within 1 second of glancing at the terminal.
- **SC-002**: The application restores the terminal 100% of the time after a user initiates a cancellation (Ctrl+C).
- **SC-003**: Enabling TUI mode does not increase the total execution time by more than 5% compared to standard text mode.