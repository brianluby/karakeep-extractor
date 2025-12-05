# CLI Command Contract: Help Text

## Global

```text
Usage: karakeep <command> [flags]

Commands:
  setup      Run interactive configuration wizard
  extract    Fetch bookmarks from Karakeep
  enrich     Enrich repositories with GitHub metadata
  rank       Display ranked list of repositories
```

## Subcommands

Ensure all flags have descriptions.

*   `extract`: `--url`, `--token`, `--db`
*   `enrich`: `--token`, `--limit`, `--force`
*   `rank`: `--limit`, `--sort`, `--tag`, `--format`, `--sink-url`, `--sink-header`, `--sink-trillium`
