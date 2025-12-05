# Quickstart: Trillium Integration

## Setup

1.  Run setup:
    ```bash
    karakeep setup
    ```
2.  Answer **Yes** when asked to configure Trillium.
3.  Provide your ETAPI URL (e.g., `http://localhost:8080`) and Token.

## Sending Data

```bash
karakeep rank --sink-trillium
```

## Verify

1.  Open Trillium.
2.  Check the Root (or Inbox).
3.  Look for a note titled "GitHub Rankings - [Date]".
