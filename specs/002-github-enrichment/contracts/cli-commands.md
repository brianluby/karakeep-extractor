# CLI Command Contract: Enrich

## Command

```bash
karakeep enrich [flags]
```

## Description

Scans the local database for GitHub URLs that have not been enriched (or all if forced), fetches metadata from GitHub API, and updates the database.

## Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--token` | `-t` | String | `$GITHUB_TOKEN` | GitHub Personal Access Token. If not provided, tries env var. If neither, runs unauthenticated. |
| `--limit` | `-l` | Int | `50` | Maximum number of repositories to process in this run. |
| `--force` | `-f` | Bool | `false` | If set, re-enriches repositories even if they already have data. |

## Output

**Success (Stdout):**
```text
Starting enrichment for 10 repositories...
[1/10] github.com/owner/repo1 ... OK (Stars: 120)
[2/10] github.com/owner/repo2 ... NOT FOUND
...
Enrichment complete.
Updated: 9
Errors: 1
Rate Limit Remaining: 4950
```

**Error (Stderr):**
```text
Error: Rate limit exceeded. Resets at 14:00.
Progress saved.
```
