# Karakeep Extractor Usage Guide

This document provides detailed instructions for using the **Karakeep Extractor** CLI tool.

## Table of Contents

1.  [Installation](#installation)
2.  [Configuration](#configuration)
3.  [Core Commands](#core-commands)
    *   [Extract](#extract)
    *   [Enrich](#enrich)
    *   [Rank](#rank)
4.  [Recipes & Workflows](#recipes--workflows)
5.  [Troubleshooting](#troubleshooting)

---

## Installation

**From Source (Go required):**

```bash
go install github.com/brianluby/karakeep-extractor/cmd/extractor@latest
# or
git clone https://github.com/brianluby/karakeep-extractor.git
cd karakeep-extractor
go build -o karakeep cmd/extractor/main.go
```

## Configuration

Run the setup wizard to configure your credentials interactively:

```bash
karakeep setup
```

This creates a configuration file at `~/.config/karakeep/config.yaml`.

### Environment Variables

You can override configuration using environment variables:

*   `KARAKEEP_URL`: Base URL of your Karakeep instance.
*   `KARAKEEP_TOKEN`: API Bearer Token for Karakeep.
*   `GITHUB_TOKEN`: Personal Access Token for GitHub API.
*   `KARAKEEP_DB`: Path to the SQLite database file.

---

## Core Commands

### Extract

Fetches bookmarks from Karakeep and saves valid GitHub repository links to the local database.

```bash
karakeep extract
```

*   **--url**: Override Karakeep URL.
*   **--token**: Override Karakeep Token.
*   **--db**: Override Database Path.

### Enrich

Fetches metadata (stars, forks, description, last pushed date) from the GitHub API for all extracted repositories.

```bash
karakeep enrich
```

*   **--limit**: Number of repositories to process (default 50).
*   **--force**: Re-process repositories even if they already have data.
*   **--token**: Override GitHub Token.

### Rank

Displays a ranked list of repositories. Supports filtering and exporting.

```bash
karakeep rank
```

*   **--sort**: Sort metric (`stars` [default], `forks`, `updated`).
*   **--limit**: Number of results to show (default 20).
*   **--tag**: Filter by keyword in Title or Description (case-insensitive).
*   **--format**: Output format (`table` [default], `json`, `csv`).
*   **--sink-url**: URL to POST the JSON result payload to.
*   **--sink-header**: Headers for sink request (Key: Value).
*   **--sink-trillium**: Enable direct export to Trillium Notes.

---

## Recipes & Workflows

### 1. Daily Review of Top Projects

Setup a quick alias to see what you've saved recently or what's popular.

```bash
# Top 10 by Stars
karakeep rank --limit 10

# Top 5 Most Recently Updated
karakeep rank --sort updated --limit 5
```

### 2. Deep Dive into a Topic

Filter your list to find tools for a specific ecosystem.

```bash
# Find all 'python' tools
karakeep rank --tag "python" --limit 50
```

### 3. Export for Spreadsheet Analysis

Dump your entire enriched database to CSV for further analysis in Excel or Google Sheets.

```bash
karakeep rank --limit 1000 --format csv > my_repos.csv
```

### 4. Sync to Trillium Notes

Keep a running "leaderboard" note in your PKM system.

```bash
# Requires 'karakeep setup' with Trillium details first
karakeep rank --sink-trillium --limit 20
```

### 5. Custom Webhook Integration

Send your top picks to a custom endpoint (e.g., Zapier, Discord webhook wrapper).

```bash
karakeep rank \
  --limit 5 \
  --sink-url "https://hooks.zapier.com/..." \
  --sink-header "Content-Type: application/json"
```

---

## Troubleshooting

### Rate Limit Exceeded

**Error**: `rate limit exceeded`

*   **Cause**: You are making too many requests to GitHub without a token (limit: 60/hr) or have exhausted your token's quota (5000/hr).
*   **Solution**:
    1.  Ensure you have configured a GitHub Token via `karakeep setup`.
    2.  Wait for the reset window (usually 1 hour).
    3.  Use `karakeep enrich --limit 10` to process in smaller batches.

### Connection Refused

**Error**: `connection refused` or `no such host`

*   **Cause**: The tool cannot reach your Karakeep or Trillium instance.
*   **Solution**:
    1.  Check if the service is running.
    2.  Verify the URL in `config.yaml` or environment variables.
    3.  Ensure you are on the correct network (VPN/LAN) if self-hosting.

### Unauthorized (401)

**Error**: `401 Unauthorized`

*   **Cause**: Invalid API Token.
*   **Solution**: Regenerate the token in Karakeep/Trillium settings and update via `karakeep setup`.
