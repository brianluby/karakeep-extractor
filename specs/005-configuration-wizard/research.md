# Research & Decision Log

## 1. Config File Format

**Context**: Need to persist user settings.
**Question**: JSON vs YAML vs TOML vs .env?

**Decision**: **YAML**

**Rationale**:
*   **Readability**: Easier for humans to read and edit manually than JSON.
*   **Standard**: Common standard for CLI tools (kubectl, helm, etc.).
*   **Library**: `gopkg.in/yaml.v3` is robust.

**Alternatives Considered**:
*   JSON: Brittle syntax (no trailing commas, strict quoting).
*   .env: Good for secrets, but less structured for potentially nested config (though our config is flat, YAML is more "config-like").

## 2. Config File Location

**Context**: Where to store the file?
**Question**: Hardcoded path vs Standard locations?

**Decision**: **XDG Config Standard (UserConfigDir)**

**Rationale**:
*   **Best Practice**: On Linux/Mac, follows `~/.config/appname`. On Windows, `%AppData%`.
*   **Implementation**: Go `os.UserConfigDir()` provides this path reliably.

## 3. Interaction Library

**Context**: Need to prompt user.
**Question**: Stdlib vs External Lib?

**Decision**: **Stdlib (`fmt`, `bufio`, `os`)**

**Rationale**:
*   **Constitution**: Minimize dependencies.
*   **Complexity**: We only need simple string prompts. No complex selection lists or colors required for MVP setup.

**Alternatives Considered**:
*   `charmbracelet/huh`: Beautiful forms, but heavy dependency for 3 questions.

## 4. Precedence Logic

**Context**: User can supply config via File, Env, or Flag.
**Question**: What overrides what?

**Decision**: **Flag > Env > File**

**Rationale**:
*   **Convention**: Explicit flags are most specific (session-scope). Env vars are environment-scope. Files are user-scope (least specific).
*   **Safety**: Allows overriding persisted secrets for a single run without editing files.
