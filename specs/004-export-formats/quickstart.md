# Quickstart: Exports & Sinks

## Prerequisites

1.  **Karakeep Database**: Ensure you have run `karakeep enrich` to populate the database.

## Exporting Data

### JSON Export

Useful for piping to `jq` or saving to disk.

```bash
karakeep rank --format json > repos.json
```

### CSV Export

Useful for spreadsheets.

```bash
karakeep rank --format csv > repos.csv
```

## Using API Sink

Send the ranked list to a webhook (e.g., RequestBin for testing).

```bash
# 1. Create a bin at requestbin.com or use a local echo server
# 2. Send data
karakeep rank --sink-url "https://your-endpoint.com/webhook"
```

### Custom Headers

If your API requires authentication:

```bash
karakeep rank \
  --sink-url "https://api.trillium.cc/notes" \
  --sink-header "Authorization: Bearer my-token" \
  --sink-header "X-Custom-Source: KarakeepCLI"
```

