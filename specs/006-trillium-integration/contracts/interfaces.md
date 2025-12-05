# Internal Interfaces: Trillium

No new domain interfaces needed (uses existing `Sink`).

However, `TrilliumClient` will expose:

```go
type TrilliumClient interface {
    CreateNote(ctx context.Context, title string, contentHTML string) error
}
```
