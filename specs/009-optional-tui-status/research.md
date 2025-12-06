# Research: Optional TUI Status

## 1. TUI Library Selection
**Decision:** Use `charmbracelet/bubbletea` (plus `bubbles` and `lipgloss`).
**Rationale:**
- **Ecosystem Standard:** It is the most popular and well-maintained TUI framework for Go.
- **Architecture:** The Elm Architecture (Model-Update-View) fits well with the state-driven nature of tracking enrichment progress.
- **Robustness:** Handles window resizing, alt-screen switching, and event loops correctly, which is hard to get right with raw `x/term` or `termbox`.
- **Components:** `charmbracelet/bubbles` provides ready-to-use spinners and progress bars, accelerating development.

**Alternatives Considered:**
- `gdamore/tcell`: Lower level, good for complex layouts but requires more boilerplate for basic components.
- `rivo/tview`: Good widget set, but less "modern" feel and harder to style custom components compared to Lipgloss.
- `pterm`: Great for simple "print and forget" output, but less capable for interactive event loops (handling Ctrl+C cleanly while updating dynamic bars).

## 2. Architecture & decoupling
**Decision:** Use the Observer Pattern via a `ProgressReporter` interface.
**Rationale:**
- The Core Services (`Enricher`, `Extractor`) must not depend on `bubbletea`.
- By passing a `ProgressReporter` interface into the service methods, we can inject a TUI-aware reporter when the `--tui` flag is present, or a no-op/standard-logger reporter when it isn't.
- This satisfies FR-007 (Legacy behavior preserved) without duplicating core logic.

**Interface Design:**
```go
type ProgressReporter interface {
    Init(total int, message string)
    ReportProgress(increment int, message string)
    ReportError(err error)
    Finish(summary string)
}
```

## 3. Handling Logs
**Decision:** Suppress standard logs in TUI mode; redirect critical errors to a "Tail" model.
**Rationale:**
- Standard `log.Printf` writes to stderr, which will corrupt the TUI rendering.
- We will create a custom implementation of `ProgressReporter` that captures these messages as "events" and feeds them into a `[]string` slice in the Bubble Tea model (The "Static Tail").
- For the legacy mode, the reporter will just pass through to `log.Printf` or `fmt.Println`.

## 4. Unknowns Resolved
- **UI-Driven Loop:** The `main` function will initialize the `tea.Program`. The core logic (`EnrichBatch` / `Extract`) will run in a `go` routine triggered by the `Init` command of the tea program.
- **Indeterminate Progress:** Since Extraction doesn't know the total, we will use a `bubbles/spinner` and a simple counter. Enrichment knows the total (from `limit` or DB count), so it will use `bubbles/progress`.