# Quickstart: using the TUI

## Prerequisites
- Terminal with ANSI support (Standard on macOS/Linux, requires modern Windows Terminal on Windows).

## Running with TUI

The `--tui` flag is available on both `extract` and `enrich` commands.

### 1. Enrichment
Visualize the enrichment progress of your extracted repositories.

```bash
# Enrich with TUI enabled
karakeep-extractor enrich --tui
```

**Expected Output:**
- A progress bar fills up across the screen.
- Current repository "owner/repo" is shown below the bar.
- Recent errors (e.g. 404s) appear at the bottom.

### 2. Extraction
Visualize the bookmark fetching process.

```bash
# Extract with TUI enabled
karakeep-extractor extract --tui
```

**Expected Output:**
- A spinner animation indicates activity.
- A counter shows "Extracted Repos: X".

### 3. Legacy Mode
If you omit the flag, the standard text logs are preserved (useful for scripts/CI).

```bash
# Standard output (no TUI)
karakeep-extractor enrich
```