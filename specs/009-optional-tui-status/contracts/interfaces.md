# Contracts: Progress Reporting

## Go Interfaces

This feature introduces internal interfaces to decouple the UI from the Core Domain.

### 1. ProgressReporter

Defined in `internal/core/domain/interfaces.go` (or similar package).

```go
package domain

// ProgressReporter abstracts the output mechanism (CLI logs vs TUI updates).
type ProgressReporter interface {
    // Start initializes the progress tracking. 
    // total: expected number of items (-1 if unknown).
    // title: description of the task.
    Start(total int, title string)

    // Increment adds to the processed count.
    Increment()

    // SetStatus updates the description of the current item being processed.
    SetStatus(status string)

    // Log records a message (info/error) without stopping the process.
    // In TUI mode, this goes to the log tail. In text mode, this is stderr/stdout.
    Log(message string)

    // Error records an error specifically (may be highlighted differently).
    Error(err error)
}
```

### 2. Service Modifications

Existing services must be updated to accept this interface.

```go
// internal/core/service/enricher.go

// EnrichBatch signature update
func (e *Enricher) EnrichBatch(ctx context.Context, limit int, force bool, workers int, reporter domain.ProgressReporter) (int, int, error)
```

```go
// internal/core/service/extractor.go

// Extract signature update
func (e *Extractor) Extract(ctx context.Context, reporter domain.ProgressReporter) error
```

### 3. TUI Program Interface

The TUI entrypoint.

```go
// internal/ui/tui/program.go

// Run starts the Bubble Tea program.
// task: A closure/function that performs the actual work and uses the internal reporter.
func Run(ctx context.Context, mode string, task func(domain.ProgressReporter) error) error
```