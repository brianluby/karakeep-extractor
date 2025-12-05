# Karakeep Extraction - Product Overview

**Version:** 2.0
**Date:** December 2025
**Status:** Concept & Planning Phase

---

## Executive Summary

This document outlines the **Karakeep Extraction** tool, a specialized command-line application designed to bridge the gap between personal knowledge collection in Karakeep and open-source project discovery on GitHub.

### Core Mission

To identify, enrich, and surface high-potential open-source projects by:
1.  **Extracting** raw project links and AI-generated summaries from a user's Karakeep instance.
2.  **Enriching** those links with real-time metadata (stars, forks, last update) from the GitHub API.
3.  **Ranking** and presenting the results to help users prioritize which projects to explore.

---

## Project Vision

### Problem Statement

Users often save interesting GitHub repositories to Karakeep (formerly Hoarder) as they browse. However, these bookmarks sit in a static list without context on their popularity or activity. Users miss out on identifying "trending" or "high-value" projects buried in their own archives.

### Solution

A CLI tool that acts as an intelligence layer over Karakeep. It doesn't just manage bookmarks; it **analyzes** them.

**Key Value Props:**
*   **Discovery:** Find the "hidden gems" in your own bookmark list.
*   **Freshness:** Know which projects are active and which are abandoned.
*   **Focus:** Prioritize reading based on community signal (stars/forks).

---

## Data Sources & Integrations

### 1. Karakeep (Source)
*   **Role:** The primary data source for "raw" project leads.
*   **Interaction:** The CLI queries the Karakeep API to fetch bookmarks, specifically looking for `github.com` URLs and their associated AI summaries/tags.
*   **Auth:** Bearer Token (JWT).

### 2. GitHub (Enrichment)
*   **Role:** The source of truth for project health and popularity.
*   **Interaction:** The CLI queries the GitHub REST API to fetch:
    *   Star Count (`stargazers_count`)
    *   Fork Count (`forks_count`)
    *   Last Updated Date (`pushed_at`)
    *   Description & Language
*   **Auth:** Personal Access Token (optional but recommended for rate limits).

---

## Phase 1: CLI Application

### Core Features

#### 1.1 Extraction & Enrichment
```bash
# The core loop: Fetch from Karakeep -> Enrich via GitHub -> Output Table
karakeep-extractor run [--limit 50] [--tag "dev-tools"]

# Example Output:
# | Rank | Project          | Stars  | Forks | Last Update | Summary (from Karakeep) |
# |------|------------------|--------|-------|-------------|-------------------------|
# | 1    | charmbracelet/   | 35k    | 1.2k  | 2 hrs ago   | TUI library for Go      |
# | 2    | karakeep/karakeep| 5k     | 200   | 1 day ago   | Bookmark manager        |
```

#### 1.2 Sync & Update
```bash
# Update metadata for existing extracted links without re-fetching from Karakeep
karakeep-extractor refresh
```

#### 1.3 Configuration
```bash
# Setup wizard for both services
karakeep-extractor setup
# Prompts for:
# - Karakeep URL & Token
# - GitHub Token (optional)
```

### Technical Architecture

#### Technology Stack
- **Language:** Golang (Go)
  - **Why:** Strong concurrency for parallel API requests (fetching GitHub stats for 100s of links).
- **HTTP Client:** `go-resty/resty`
- **CLI Framework:** `spf13/cobra`
- **Output:** `charmbracelet/lipgloss` (styling) & `olekukonko/tablewriter`.

#### Data Flow
1.  **User** runs `extract` command.
2.  **App** calls Karakeep API -> Returns List of Bookmarks.
3.  **App** filters for `github.com/*` URLs.
4.  **App** spawns concurrent workers to call GitHub API for each repo.
5.  **App** aggregates results, calculates a "Score" (default: star count), and sorts.
6.  **App** renders the ranked table to stdout.

#### Project Structure
```
karakeep-extractor/
├── cmd/
│   └── extractor/     # Main entry point
│       └── main.go
├── internal/
│   ├── adapter/       # API Clients
│   │   ├── karakeep/
│   │   └── github/
│   ├── core/          # Business Logic
│   │   ├── domain/    # Structs (Bookmark, RepoStats)
│   │   └── service/   # Orchestrator (Extract -> Enrich -> Rank)
│   ├── config/        # Viper setup
│   └── ui/            # TUI rendering
├── go.mod
└── README.md
```

---

## Development Roadmap

### Sprint 1: Foundation (Weeks 1-2)
- [ ] Project scaffolding (Go mod, Cobra, Viper).
- [ ] **Karakeep Client:** Implement `GetBookmarks` with pagination.
- [ ] **GitHub Client:** Implement `GetRepoDetails` with rate limit handling.

### Sprint 2: Core Logic (Weeks 3-4)
- [ ] **Extraction Service:** Regex parsing of GitHub URLs from bookmark links/text.
- [ ] **Enrichment Service:** Worker pool pattern to fetch GitHub stats concurrently.
- [ ] **Ranking Engine:** Basic sorting by Star Count.
- [ ] **CLI Output:** Formatted table display.

### Sprint 3: Polish (Week 5)
- [ ] Caching (don't re-fetch GitHub stats if < 24h old).
- [ ] JSON/CSV export flags (`--format=json`).
- [ ] Comprehensive Error Handling (invalid tokens, API down).

---

## Success Metrics
- **Accuracy:** Correctly identifies valid GitHub links from mixed bookmark data.
- **Performance:** Can enrich 100 bookmarks in < 5 seconds (using concurrency).
- **Utility:** User can successfully identify their "top rated" saved project.

---

## Risks & Mitigation
| Risk | Impact | Mitigation |
|------|--------|------------|
| GitHub Rate Limits | High | Implement authenticated requests (5000 req/hr) & caching. |
| Karakeep API Changes | Medium | Strict interface decoupling. |
| Data Noise | Medium | Robust regex filters to ignore non-project GitHub links (e.g., issues, profiles). |

---

**Document Version:** 2.0
**Last Updated:** December 4, 2025
**Status:** Draft - Ready for Review
