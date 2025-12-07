# Contracts: Real-time TUI Statistics

## ProgressReporter Interface

The core contract for reporting progress is updated.

```go
package domain

type ProgressReporter interface {
	Start(total int, title string)
	Increment()
	SetStatus(status string)
	Log(message string)
	Error(err error)
	Finish(summary string)

    // New additions
    RecordSuccess()
    RecordFailure()
    RecordSkipped()
}
```

### Behavior in `TextReporter` (CLI mode)
- `RecordSuccess()`: No output (too verbose), or maybe a dot `.`?
- `RecordFailure()`: Log error to stderr.
- `RecordSkipped()`: Log "Skipped" to stdout (verbose mode only?).

### Behavior in `BubbleTeaReporter` (TUI mode)
- `RecordSuccess()`: Send `MsgSuccess`.
- `RecordFailure()`: Send `MsgFailure`.
- `RecordSkipped()`: Send `MsgSkipped`.
