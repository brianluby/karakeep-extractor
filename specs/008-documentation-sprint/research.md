# Research & Decision Log

## 1. Documentation Structure

**Context**: How to organize the documentation for best usability?
**Question**: Single massive README vs Split Docs?

**Decision**: **Hybrid (README + docs/)**

**Rationale**:
*   **README**: Must contain "Quick Start", "Installation", and "Basic Usage". It's the landing page.
*   **docs/usage.md**: Detailed flag reference, troubleshooting, and complex recipes. Keeps the README clean but provides depth for power users.

## 2. Help Text Sync

**Context**: `karakeep --help` should match the docs.
**Question**: Auto-generate docs from help or manual sync?

**Decision**: **Manual Sync (for now)**

**Rationale**:
*   **Effort**: Building a doc generator is overkill for 4 commands.
*   **Quality**: Hand-written docs are often clearer than auto-generated flag lists.
*   **Action**: We will manually verify the help strings in `main.go` during this sprint.
