# Quickstart: GitHub Repository Ranking

## Prerequisites

1.  **Karakeep Database**: Ensure you have run `karakeep enrich` to populate the database with stars/forks.
    *   *If you haven't enriched yet, the ranking will return no results or 0s.*

## Running Ranking

### 1. Top 20 by Stars (Default)

```bash
karakeep rank
```

### 2. Top 5 Most Recent

```bash
karakeep rank --limit 5 --sort updated
```

### 3. Top 50 by Forks

```bash
karakeep rank --limit 50 --sort forks
```

## Expected Output

You should see a clean, tab-aligned table in your terminal. If the list is long, it may open in `less` automatically (press `q` to exit).
