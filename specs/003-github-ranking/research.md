# Research & Decision Log

## 1. Table Rendering Library

**Context**: We need to display a formatted ASCII table of repositories.
**Question**: Should we build a custom table printer or use an existing library?

**Decision**: **Use `olekukonko/tablewriter` (or similar lightweight lib)**

**Rationale**:
*   **Maintainability**: Correctly handling column alignment, padding, and unicode characters (for potential emojis in descriptions) is tedious to implement from scratch.
*   **Features**: Libraries often support auto-wrapping, headers, and borders out of the box.
*   **Constitution**: "Minimize external dependencies" is key, but for UI formatting, a small, stable library is often better than buggy custom code.
*   *Alternative*: Standard `text/tabwriter`.

**Revised Decision**: **Use `text/tabwriter` (Standard Library)**

**Revised Rationale**:
*   **Zero Dependency**: It's built-in.
*   **Sufficiency**: For a simple rank list, we don't need complex borders. A clean tab-separated list is often more CLI-native and pipe-friendly anyway.
*   **Formatting**: `tabwriter` handles column alignment perfectly. We can format headers manually.

**Alternatives Considered**:
*   `olekukonko/tablewriter`: Rich features but external dep.
*   `charmbracelet/lipgloss`: Beautiful but overkill/heavy dep.

## 2. Sorting Implementation

**Context**: Users can sort by `stars`, `forks`, or `updated`.
**Question**: Should sorting happen in the Database (SQL) or in Application Memory (Go)?

**Decision**: **Database (SQL)**

**Rationale**:
*   **Performance**: SQLite is optimized for sorting. `ORDER BY stars DESC LIMIT 20` is extremely fast and avoids fetching thousands of rows into memory just to show the top 10.
*   **Scalability**: If the dataset grows, in-memory sorting becomes a bottleneck.
*   **Simplicity**: The query `SELECT ... WHERE enrichment_status='SUCCESS' ORDER BY ... LIMIT ...` is clean and declarative.

**Alternatives Considered**:
*   In-memory sort: More flexible for complex custom logic, but slower and resource-intensive for large datasets.

## 3. Paged Output

**Context**: Output might exceed terminal height.
**Question**: How to implement paging?

**Decision**: **Detect TTY and pipe to `less` command**

**Rationale**:
*   **UX**: Standard unix behavior.
*   **Implementation**: Go can detect if stdout is a TTY. If so, and row count is high, we can execute `less` and pipe our output to its stdin.
*   **Fallback**: If `less` is missing or not TTY (piped to file), just print.

**Alternatives Considered**:
*   Internal pager (implementing "Press Enter for more"): Tedious and reinvents the wheel.
