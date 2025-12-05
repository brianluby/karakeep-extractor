# Data Model: Trillium Integration

## Configuration

New fields in `config.yaml`:

```yaml
trillium_url: "http://trillium.local:8080"
trillium_token: "etapi-token-..."
```

## API Payloads

### Create Note Request

Target: `POST /etapi/create-note`

```json
{
  "parentNoteId": "root", 
  "title": "GitHub Rankings - 2023-10-27",
  "type": "text",
  "content": "<table>...</table>" 
}
```
*Note: `parentNoteId` could be configurable, but "root" or "inbox" (if search finds it) is a safe default. For MVP, let's target "root".*

## Internal Structures

### TrilliumSink

Implements `domain.Sink`.

```go
type TrilliumSink struct {
    client *TrilliumClient
}
```
