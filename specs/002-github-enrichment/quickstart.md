# Quickstart: GitHub Enrichment

## Prerequisites

1.  **Karakeep Database**: Ensure you have run `karakeep extract` previously to populate the database with URLs.
2.  **GitHub Token (Optional but Recommended)**:
    *   Generate a Personal Access Token (Classic) with `public_repo` scope (or no scope for public-only read).
    *   Export it: `export GITHUB_TOKEN=ghp_xxxxxxxxxxxx`

## Running Enrichment

### 1. Standard Run

Enrich up to 50 pending repositories using the environment token.

```bash
karakeep enrich
```

### 2. Custom Limit & Token

Enrich only 10 repositories using a specific token passed via flag.

```bash
karakeep enrich --limit 10 --token ghp_mytoken123
```

### 3. Force Re-Enrichment

Update data for repositories that were already processed.

```bash
karakeep enrich --force
```

## Verifying Results

Check the SQLite database directly:

```bash
sqlite3 karakeep.db "SELECT url, stars, enrichment_status FROM extracted_repos WHERE stars IS NOT NULL LIMIT 5;"
```
