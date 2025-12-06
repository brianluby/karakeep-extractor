# Usage Guide

## Basic Commands

### Extraction

Fetch bookmarks from Karakeep and save potential GitHub repositories to your local database.

```bash
# Standard extraction (text logs)
karakeep-extractor extract

# Extraction with Visual Status (TUI)
karakeep-extractor extract --tui
```

### Enrichment

Fetch metadata (stars, forks, description) from GitHub for the repositories you have extracted.

```bash
# Standard enrichment (text logs)
karakeep-extractor enrich

# Enrichment with Visual Progress (TUI)
karakeep-extractor enrich --tui

# Options
karakeep-extractor enrich --limit 100  # Process up to 100 repos
karakeep-extractor enrich --force      # Re-process already enriched repos
```

### Ranking

View your top repositories.

```bash
# Top 20 by stars (default)
karakeep-extractor rank

# Top 50 by forks
karakeep-extractor rank --sort forks --limit 50

# Export to CSV
karakeep-extractor rank --format csv > ranking.csv
```

### Setup

Configure your API tokens interactively.

```bash
karakeep-extractor setup
```