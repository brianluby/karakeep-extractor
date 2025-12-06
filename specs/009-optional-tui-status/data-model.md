# Data Model: TUI State

## TUI State Models

Since this feature is UI-centric, the "Data Model" refers to the application state held by the Bubble Tea program.

### 1. Root Model (`TuiModel`)

The top-level container for the TUI application.

| Field | Type | Description |
| :--- | :--- | :--- |
| `State` | `AppState` | Enum: `Idle`, `Running`, `Done`, `Error` |
| `Mode` | `OperationMode` | Enum: `Extract`, `Enrich` |
| `Progress` | `progress.Model` | Bubble Tea Progress Bar component (for Enrichment) |
| `Spinner` | `spinner.Model` | Bubble Tea Spinner component (for Extraction) |
| `Count` | `int` | Number of items processed so far |
| `Total` | `int` | Total items to process (Enrichment only, -1 for Extract) |
| `CurrentItem` | `string` | Name/ID of the item currently being processed |
| `Logs` | `[]string` | Ring buffer of last N log messages/errors |
| `Err` | `error` | Fatal error if process crashed |

### 2. Events (Messages)

These structure the data flow from the Worker (Core Logic) to the UI.

| Message Type | Payload | Description |
| :--- | :--- | :--- |
| `MsgStart` | `Total int`, `Title string` | Signals the start of the batch operation. |
| `MsgProgress` | `ItemName string`, `Inc int` | Updates progress. `Inc` usually 1. |
| `MsgLog` | `Level string`, `Text string` | Adds a line to the log tail. |
| `MsgDone` | `Summary string` | Signals completion. |
| `MsgFatal` | `Err error` | Signals a crash/stop. |

### 3. Configuration

No persistent database changes required for this feature.
Config is runtime-only via CLI flags.