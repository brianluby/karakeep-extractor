# Internal Interfaces: Export & Sink

## Exporter

Responsible for formatting the output to a specific stream (usually stdout).

```go
package domain

import "io"

type ExportFormat string

const (
    FormatTable ExportFormat = "table"
    FormatJSON  ExportFormat = "json"
    FormatCSV   ExportFormat = "csv"
)

type Exporter interface {
    // Export writes the repositories to the provided writer in the specific format.
    Export(repos []ExtractedRepo, w io.Writer) error
}
```

## Sink

Responsible for sending the data to an external system.

```go
package domain

import "context"

type Sink interface {
    // Send transmits the repository list to the configured endpoint.
    // It assumes the payload is always JSON.
    Send(ctx context.Context, repos []ExtractedRepo) error
}
```
