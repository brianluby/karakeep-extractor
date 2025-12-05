# CLI Command Contract: Filtering

## Command

```bash
karakeep rank --tag "keyword"
```

## Description

Filters the ranked list to show only repositories containing "keyword" in their title or description.

## Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--tag` | String | `""` | Keyword to filter repositories by (case-insensitive). |

## Example

```bash
$ karakeep rank --tag "cli"
RANK  NAME              STARS  FORKS  UPDATED
1     charmbracelet/bub 12000  500    2d ago
...
```
