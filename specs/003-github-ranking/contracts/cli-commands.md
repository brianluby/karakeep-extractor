# CLI Command Contract: Rank

## Command

```bash
karakeep rank [flags]
```

## Description

Displays a ranked list of extracted GitHub repositories.

## Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--limit` | `-l` | Int | `20` | Number of repositories to display. |
| `--sort` | `-s` | String | `stars` | Metric to sort by. Values: `stars`, `forks`, `updated`. |

## Output (Stdout)

**Header:**
`RANK  NAME              STARS  FORKS  UPDATED`

**Rows (Tab-aligned):**
```text
1     charmbracelet/bub 12000  500    2d ago
2     junegunn/fzf      9000   300    1w ago
...
```

## Output (Stderr)

*   Error messages (e.g., "Invalid sort parameter").
