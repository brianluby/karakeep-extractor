# Feature Roadmap

This document outlines potential new features and enhancements for the Karakeep Extractor.

## Planned Features

### 1. Real-time TUI Statistics
**Description:** Enhance the TUI (Text User Interface) to display real-time counters during extraction and enrichment processes.
**Goal:** Provide immediate visual feedback on progress, showing exactly how many items have been successfully processed versus failed or skipped, updating dynamically below the progress bar.

### 2. LLM-Powered Repository Analysis
**Description:** Integrate with an LLM (Large Language Model) API to perform advanced analysis on the extracted repository data.
**Functionality:** Select data points like description, language, stars, and forks from the local database and send them to a configurable LLM.
**Use Case:** Enable complex queries such as "Which of these Go projects would be most useful for building a distributed system?" or "Summarize the trending Python tools from my bookmarks."

### 3. Karakeep AI Tags Extraction
**Description:** Extract AI-generated tags associated with bookmarks from the Karakeep API.
**Requirement:** Research the Karakeep API response structure to identify where these tags are stored and map them to the local database schema. This will improve filtering and categorization capabilities.

### 4. Advanced Local Database Querying
**Description:** Add a powerful CLI command to query the local SQLite database with flexible filters.
**Capabilities:** Support filtering by multiple keywords, star count ranges (e.g., `>1000`), fork counts, dates (e.g., `created_after=2024-01-01`), and combinations of these criteria.

## Future Enhancements (Ideas)

### 5. Automatic Periodic Sync
**Description:** Implement a background service or cron-friendly mode to automatically fetch new bookmarks and enrich them on a schedule, keeping the local database up-to-date without manual intervention.

### 6. "Trending" & "Stale" Reports
**Description:** Generate reports highlighting repositories that have gained significant traction (stars/forks) since they were bookmarked, or identifying "stale" projects that haven't been updated in years.

### 7. Interactive Terminal Browser
**Description:** Develop a fully interactive TUI mode that allows users to browse, search, sort, and filter their extracted repositories directly within the terminal, similar to a file manager but for code repos.

### 8. Export to Markdown / Obsidian / Notion
**Description:** Create specialized export formats optimized for personal knowledge management tools. For example, generating a single Markdown file with a categorized list of repositories, or individual files for each repo with metadata as frontmatter.
