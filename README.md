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

##  roadmap

- [ ] **v1 (CLI):** Core logic for extraction, enrichment, and outputting ranked lists to the terminal or JSON/CSV.
- [ ] **v2 (Web UI):** A simple dashboard to visualize the ranked projects and filter by tags/categories.

## 🛠️ Architecture

1.  **Source:** KaraKeep API
2.  **Processor:** Parses content, extracts `github.com` links.
3.  **Enricher:** Async queries to GitHub API.
4.  **Output:** Ranked list.

---
*Created for the Karakeep Extractor project.*
