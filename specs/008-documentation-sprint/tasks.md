# Tasks: Documentation Sprint

**Feature**: `008-documentation-sprint`
**Spec**: [specs/008-documentation-sprint/spec.md](spec.md)

## Implementation Strategy
- **Phase 1: Setup**: Create docs folder.
- **Phase 2: Usage Guide**: Write detailed command reference.
- **Phase 3: README**: Rewrite root README to point to new docs and cover basics.
- **Phase 4: CLI Help**: Align code help strings with documentation.
- **Phase 5: Polish**: Review and format.

---

## Phase 1: Setup

- [x] T001 Create `docs/` directory

---

## Phase 2: Usage Guide (User Story 2)

**Goal**: Provide deep-dive documentation.
**Priority**: P2

- [x] T002 Write `docs/usage.md` with Command Reference (Setup, Extract, Enrich, Rank)
- [x] T003 Add "Recipes" section to `docs/usage.md` (Daily Review, CSV Export, Trillium Sync)
- [x] T004 Add "Troubleshooting" section to `docs/usage.md`

---

## Phase 3: README Rewrite (User Story 1)

**Goal**: Create a clear landing page.
**Priority**: P1

- [x] T005 Rewrite `README.md` to include: Introduction, Features, Installation, Quickstart, and link to `docs/usage.md`

---

## Phase 4: CLI Help Sync (User Story 1)

**Goal**: Ensure inline help is helpful.
**Priority**: P1

- [x] T006 Update usage strings in `cmd/extractor/main.go` for `extract`, `enrich`, `rank`, `setup` commands to match new documentation clarity

---

## Phase 5: Polish

- [x] T007 Verify markdown rendering
- [x] T008 Verify `karakeep --help` output

---

## Dependencies

1. **T001** must complete before **T002**.
2. **T002-T004** should ideally be drafted before finalizing **T005** (README links).

## Parallel Execution Examples

- **Team A**: Write Markdown Docs (T002-T005).
- **Team B**: Update Go Code Help Strings (T006).
