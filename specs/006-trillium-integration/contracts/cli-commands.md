# CLI Command Contract: Trillium

## Command

```bash
karakeep rank --sink-trillium
```

## Description

Generates the ranked list and pushes it to the configured Trillium instance.

## Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--sink-trillium` | Bool | `false` | Enable Trillium sink. |

*Note: Requires configuration via `karakeep setup`.*

## Setup Interaction

```text
Configure Trillium Integration? [y/N]: y
Enter Trillium Instance URL: http://my-trilium.com
Enter Trillium ETAPI Token: [hidden]
```
