# Karakeep Extractor

**Karakeep Extractor** is a tool designed to surface high-potential open-source projects by bridging data from Karakeep with metadata from GitHub.

## 🎯 Objective

The primary goal is to identify projects of interest by:
1.  **Extracting** GitHub URLs and AI summaries from KaraKeep posts.
2.  **Enriching** this data by connecting to the GitHub API to retrieve star counts and other metadata.
3.  **Ranking** the results to prioritize the most popular or trending repositories.

## ✨ Features

- **Data Extraction:** Scrapes/Queries KaraKeep API for post summaries and embedded links.
- **GitHub Integration:** Fetches real-time repository statistics (Stars, Forks, Last Updated).
- **Force Ranking:** Sorts projects based on popularity metrics to aid in discovery.
- **CLI Interface:** Simple command-line tool for execution.

## 🚀 Quickstart

### Prerequisites
*   Go 1.25+
*   Access to a Karakeep instance (URL + Token)

### Build
```bash
go build -o karakeep-extractor cmd/extractor/main.go
```

### Run Extraction

1.  **Set Environment Variables** (Optional but recommended):
    ```bash
    export KARAKEEP_URL="https://my-karakeep.com"
    export KARAKEEP_TOKEN="my-secret-token"
    ```

2.  **Run the command**:
    ```bash
    ./karakeep-extractor extract
    ```

3.  **Verify Output**:
    ```bash
    sqlite3 karakeep.db "SELECT * FROM extracted_repos;"
    ```

## 🔌 Integrations & Exports

Future versions will expand the CLI's capability to pipe data into personal knowledge management systems and other APIs.

- **JSON/CSV Export**: Standardize output for processing by other tools (`--format=json`).
- **Trillium Notes**: Direct API integration to create notes for top-ranked repositories.
- **Webhook Stubs**: Generic HTTP POST capability to send ranking payloads to user-defined endpoints (e.g., Zapier, n8n).

## 🗺️ Roadmap

- [ ] **v1 (CLI):** Core logic for extraction, enrichment, and outputting ranked lists to the terminal or JSON/CSV.
- [ ] **v2 (Web UI):** A simple dashboard to visualize the ranked projects and filter by tags/categories.

## 🛠️ Architecture

1.  **Source:** KaraKeep API
2.  **Processor:** Parses content, extracts `github.com` links.
3.  **Enricher:** Async queries to GitHub API.
4.  **Output:** Ranked list.

---
*Created for the Karakeep Extractor project.*