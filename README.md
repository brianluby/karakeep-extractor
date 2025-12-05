# Karakeep Extractor

[![Go Report Card](https://goreportcard.com/badge/github.com/brianluby/karakeep-extractor)](https://goreportcard.com/report/github.com/brianluby/karakeep-extractor)
[![License](https://img.shields.io/github/license/brianluby/karakeep-extractor)](LICENSE)

**Karakeep Extractor** is a powerful CLI tool designed to surface high-potential open-source projects by bridging your personal knowledge base (Karakeep) with real-time metadata from GitHub.

It helps you turn a static list of bookmarks into an actionable, ranked intelligence report.

## ✨ Key Features

- **📥 Extract**: Automatically fetch bookmarks from your Karakeep instance and identify GitHub repositories.
- **⚡ Enrich**: Fetch real-time statistics (Stars, Forks, Last Updated) from the GitHub API.
- **🏆 Rank**: Sort repositories by popularity or freshness to prioritize your reading list.
- **🔍 Filter**: Slice your data by keywords (tags) to focus on specific topics (e.g., "python", "cli").
- **📤 Export**: Output data to JSON, CSV, or pipe it directly to external APIs (like Trillium Notes).

## 📦 Installation

### From Source (Go 1.25+)

```bash
git clone https://github.com/brianluby/karakeep-extractor.git
cd karakeep-extractor
go build -o karakeep cmd/extractor/main.go
# Move to your PATH
sudo mv karakeep /usr/local/bin/
```

## 🚀 Quick Start

### 1. Setup
Configure your API credentials once. This saves them to `~/.config/karakeep/config.yaml`.

```bash
karakeep setup
```

### 2. Extract
Fetch your latest bookmarks.

```bash
karakeep extract
```

### 3. Enrich
Fetch metadata for the extracted repositories (respects rate limits).

```bash
karakeep enrich
```

### 4. Rank
View your top 20 repositories by star count.

```bash
karakeep rank
```

## 📚 Documentation

For detailed usage instructions, command references, and advanced recipes (exporting, webhooks), please see the **[Usage Guide](docs/usage.md)**.

## 🗺️ Roadmap

- [x] **v1 (CLI):** Core logic for extraction, enrichment, and ranking.
- [x] **Exporting:** JSON/CSV and Webhook support.
- [x] **Integrations:** Trillium Notes support.
- [ ] **v2 (Web UI):** A simple dashboard to visualize the ranked projects.

---
*Created for the Karakeep Extractor project.*
