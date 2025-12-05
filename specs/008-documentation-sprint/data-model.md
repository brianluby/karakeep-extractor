# Data Model: Documentation

## Content Inventory

### README.md Structure
1.  **Project Title & Badges**
2.  **Introduction**: "What is this?" (Elevator pitch)
3.  **Key Features**: Bullet points (Extract, Enrich, Rank, Export).
4.  **Installation**:
    *   Go Install
    *   Build from Source
5.  **Quick Start**:
    *   Setup
    *   Extract
    *   Rank
6.  **Roadmap**: What's next?

### docs/usage.md Structure
1.  **Command Reference**:
    *   `setup`
    *   `extract`
    *   `enrich`
    *   `rank`
2.  **Workflows (Recipes)**:
    *   "Daily Review": Setup -> Extract -> Rank (Top 10)
    *   "Deep Dive": Extract -> Enrich (Force) -> Rank (CSV) -> Spreadsheet
    *   "Automated Archiving": Cron job -> Rank (JSON) -> Trillium
3.  **Troubleshooting**:
    *   Common errors (401, Rate Limits).
