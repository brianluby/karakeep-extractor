# Feature Specification: Real-time TUI Statistics

**Feature Branch**: `011-real-time-tui-stats`  
**Created**: 2025-12-06  
**Status**: Draft  
**Input**: User description: "Real-time TUI Statistics **Description:** Enhance the TUI (Text User Interface) to display real-time counters during extraction and enrichment processes. **Goal:** Provide immediate visual feedback on progress, showing exactly how many items have been successfully processed versus failed or skipped, updating dynamically below the progress bar."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Real-time Enrichment Feedback (Priority: P1)

As a user running the enrichment process (which can take a long time due to API rate limits), I want to see exactly how many repositories have been successfully enriched, how many failed, and how many were skipped, updating in real-time.

**Why this priority**: Enrichment is the longest-running process. Blindly waiting without knowing if errors are accumulating is a poor user experience.

**Independent Test**: Can be tested by running `karakeep enrich --tui` on a database with a mix of valid and invalid repositories and observing the counters incrementing.

**Acceptance Scenarios**:

1. **Given** the enrichment process is running, **When** a repository is successfully enriched, **Then** the "Success" counter increments by 1 immediately.
2. **Given** the enrichment process is running, **When** a repository fails (e.g., 404), **Then** the "Failed" counter increments by 1 immediately.
3. **Given** the enrichment process is running, **When** a repository is skipped (e.g., already fresh), **Then** the "Skipped" counter increments by 1 immediately.
4. **Given** the process completes, **Then** the final counts displayed match the actual database results.

---

### User Story 2 - Real-time Extraction Feedback (Priority: P2)

As a user running the extraction process, I want to see counters for processed bookmarks to verify that data is actually being pulled.

**Why this priority**: Extraction is usually faster, but visibility is still important for large datasets.

**Independent Test**: Can be tested by running `karakeep extract --tui` and observing counters.

**Acceptance Scenarios**:

1. **Given** the extraction process is running, **When** a bookmark is processed, **Then** the "Processed" counter updates.
2. **Given** invalid data is encountered, **Then** the "Failed" or "Ignored" counter updates.

---

### Edge Cases

- **Terminal Resize**: How does the layout handle narrow terminals? (Should wrap or hide less critical info, but counters are critical).
- **Rapid Updates**: If processing is very fast (e.g., local skipping), the UI should not flicker or lag.
- **Zero Items**: If no items are found, counters should display 0 cleanly.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The TUI MUST display three distinct counters during processing: "Success", "Failed", and "Skipped" (or equivalent context-aware labels).
- **FR-002**: The counters MUST update immediately upon the completion of a single item's processing.
- **FR-003**: The counters MUST be visible simultaneously with the progress bar.
- **FR-004**: The TUI layout MUST accommodate at least 4-digit numbers for each counter without breaking alignment.
- **FR-005**: The final state of the counters MUST remain visible after the process finishes (until the user exits).

### Key Entities *(include if feature involves data)*

- **ProgressStats**: A transient data structure holding `SuccessCount`, `FailureCount`, `SkippedCount`.
- **TUIModel**: The Bubble Tea model state, updated to include `ProgressStats`.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: UI updates occur within 100ms of the underlying event.
- **SC-002**: Users can accurately report the number of failures without checking logs.
- **SC-003**: The layout remains stable (no visual jitter) even when processing >50 items per second.