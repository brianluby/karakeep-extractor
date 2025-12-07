# Research: Advanced Local Database Querying

**Status**: Complete
**Date**: 2025-12-06

## Decisions

### 1. Query Interface Design
**Decision**: Implement a new `query` subcommand (or extend `rank` as `query` alias) with powerful flags for filtering.
**Rationale**:
- `rank` is currently for listing "top" items. `query` implies searching for specific items. Separation is cleaner but overlapping. We will create a new `query` command to support complex filters like date ranges and numeric comparisons which `rank` does not support.

### 2. Filter Syntax
**Decision**: Use standard CLI flags with string values that support operators for flexibility.
- `--stars ">1000"`
- `--forks "100..200"`
- `--after "2024-01-01"`
- `--before "2024-01-01"`
- `--lang "Go"`
**Rationale**:
- String flags with internal parsing allow for natural expressions like ">1000" without complex CLI flag structures (like `--min-stars`, `--max-stars` for every metric).

### 3. SQL Generation Strategy
**Decision**: Dynamically build the `WHERE` clause based on parsed filters.
**Rationale**:
- SQLite does not have a sophisticated query builder wrapper in this project. String concatenation with parameterized arguments is safe and effective for this scale.

### 4. Date Handling
**Decision**: Support ISO 8601 dates (`YYYY-MM-DD`) for input and query against `found_at` (or `last_pushed_at` if specified? Spec says "default to FoundAt").
**Rationale**:
- ISO 8601 is standard and sorts lexically string-wise (though SQLite stores datetimes as strings often, native comparators work if format is consistent).

## Alternatives Considered

### Extending `rank` command
- **Pros**: Fewer commands.
- **Cons**: `rank` is about *sorting* primarily. `query` is about *filtering*. Mixing ranges into `rank` might clutter help.
- **Result**: Rejected. New `query` command is clearer.

### Using a DSL (Domain Specific Language)
- **Idea**: `karakeep query "stars > 1000 AND lang = Go"`
- **Pros**: Very flexible.
- **Cons**: Requires writing a lexer/parser. Overkill.
- **Result**: Rejected. Flags are easier to implement and use from shell.

## Open Questions Resolved

- **Date Field**: Default to `found_at` (when user extracted it) or `last_pushed_at` (activity)? Spec FR-004 says "filtering on `FoundAt` by default".
    - *Action*: Implement filtering on `found_at` initially.
