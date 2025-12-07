# Research: Real-time TUI Statistics

**Status**: Complete
**Date**: 2025-12-06

## Decisions

### 1. Stats State Management
**Decision**: Create a `ProgressStats` struct and embed it within `EnrichModel` and `ExtractModel` (or their shared state).
**Rationale**:
- **Separation of Concerns**: Keeps the statistics logic separate from the view rendering logic.
- **Flexibility**: Can be easily expanded to include other metrics in the future.
- **Simplicity**: No need for a complex global state management system; just pass messages to update the struct.

### 2. Message Passing
**Decision**: Introduce specific messages for incrementing each counter: `MsgSuccess`, `MsgFailure`, `MsgSkipped`.
**Rationale**:
- **Granularity**: Allows the `Reporter` to signal specific outcomes, which the model can then aggregate.
- **Clarity**: Makes the update loop explicitly handle each case rather than inferring from a generic log message.

### 3. Layout Strategy
**Decision**: Render statistics below the progress bar but above the logs. Use a simple horizontal layout (e.g., `Success: 10 | Failed: 2 | Skipped: 5`).
**Rationale**:
- **Visibility**: Keeps critical info near the progress indicator.
- **Responsiveness**: Horizontal layout adapts well to standard terminal widths.

## Alternatives Considered

### Parsing Log Messages
- **Idea**: Have the model regex-parse `MsgLog` content to update counters.
- **Pros**: No API change for `Reporter` interface.
- **Cons**: Brittle, implicit coupling between log text and UI state.
- **Result**: Rejected. Explicit messages are safer.

### Global Stats Singleton
- **Idea**: Use a global variable for stats.
- **Pros**: Easy access.
- **Cons**: Testing nightmare, race conditions (though Bubble Tea is sequential, the input to `Send` is concurrent).
- **Result**: Rejected. Keep state localized to the model.

## Open Questions Resolved

- **Reporter Interface**: We need to extend `domain.ProgressReporter` (or add a specialized interface) to support these fine-grained events. Alternatively, we can add `RecordSuccess()`, `RecordFailure()`, `RecordSkipped()` to the interface.
    - **Decision**: Update `domain.ProgressReporter` with `RecordSuccess`, `RecordFailure`, `RecordSkipped` to support this feature natively. This is a minor breaking change for the interface but improves type safety across the board.
    - **Backwards Compatibility**: Existing implementations (like `TextReporter`) will just log these events to stdout/stderr.
