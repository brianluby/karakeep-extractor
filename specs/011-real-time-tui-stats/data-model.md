# Data Model: Real-time TUI Statistics

## UI Entities

### ProgressStats
Used within the TUI models to track real-time progress.

```go
type ProgressStats struct {
    SuccessCount int
    FailureCount int
    SkippedCount int
}
```

## Interface Updates

### domain.ProgressReporter
Updated to support granular status reporting.

```go
type ProgressReporter interface {
    // Existing methods
    Start(total int, title string)
    Increment() // Generic progress (can remain for "processing started")
    SetStatus(status string)
    Log(message string)
    Error(err error)
    Finish(summary string)

    // New methods for statistics
    RecordSuccess()
    RecordFailure()
    RecordSkipped()
}
```

## Message Types (Bubble Tea)

```go
type MsgSuccess struct{}
type MsgFailure struct{}
type MsgSkipped struct{}
```
