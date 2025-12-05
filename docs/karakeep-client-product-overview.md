# Karakeep Client Application - Product Overview

**Version:** 1.1
**Date:** December 2025
**Status:** Concept & Planning Phase

---

## Executive Summary

This document outlines a new client application for Karakeep, a self-hosted bookmark management system with AI-powered tagging. The application will begin as a command-line interface (CLI) tool and evolve into a web-based user interface, providing users with flexible ways to interact with their Karakeep instance.

### What is Karakeep?

Karakeep (formerly Hoarder) is an open-source, self-hostable "bookmark everything" application that allows users to:
- Save web links with automatic metadata extraction
- Store text notes and snippets
- Archive images and PDF documents
- Automatically tag content using AI (OpenAI or local Ollama models)
- Perform full-text searches across all saved content
- Organize bookmarks into lists and with tags
- Create highlights and annotations
- Protect against link rot with full page archival

---

## Project Vision

### Problem Statement

While Karakeep provides excellent web and mobile interfaces, there are scenarios where users would benefit from:
1. **Quick CLI access** for power users and automation scripts
2. **Custom workflows** not available in the main interface
3. **Batch operations** on bookmarks
4. **Integration with other command-line tools** (pipe data in/out)
5. **Alternative UI experiences** tailored to specific use cases

### Solution

Build a progressive application that:
- **Phase 1:** Starts as a powerful CLI tool for Karakeep power users
- **Phase 2:** Evolves into a specialized web interface with unique features
- **Future:** Could expand to include browser extensions, desktop apps, or integrations

---

## Karakeep API Overview

### Authentication
- **Method:** JWT Bearer Token authentication
- **Security:** HTTP Bearer scheme
- Tokens are generated from the Karakeep web interface under user settings

### Available API Endpoints

Based on the Karakeep API documentation (v0.29.0), the following resources are available:

#### 1. **Bookmarks API**
- Get all bookmarks (with pagination and filtering)
- Create new bookmarks (links, notes, assets)
- Update bookmark details
- Delete bookmarks
- Manage bookmark metadata (title, description, tags)

#### 2. **Lists API**
- Get all bookmark lists
- Create custom lists
- Add/remove bookmarks from lists
- Organize bookmarks hierarchically

#### 3. **Tags API**
- Retrieve all tags
- Create custom tags
- Attach/detach tags from bookmarks
- View tag usage statistics

#### 4. **Highlights API**
- Get highlights from saved content
- Create new highlights
- Annotate bookmarked content

#### 5. **Users API**
- Get current user information
- Manage user preferences

#### 6. **Assets API**
- Upload images and PDFs
- Manage asset metadata

#### 7. **Admin API** (if admin privileges)
- Update user information
- Manage system-wide settings

#### 8. **Backups API**
- Get backup information
- Trigger backup operations

---

## Phase 1: CLI Application

### Core Features

#### 1.1 Basic Operations
```bash
# Add a bookmark
karakeep add <url> [--tags tag1,tag2] [--list "My List"]
karakeep note "Quick thought" [--tags idea]
karakeep upload image.png [--tags screenshot]

# Search bookmarks
karakeep search "machine learning"
karakeep search --tag python --list "Dev Resources"

# List operations
karakeep list                    # Show all lists
karakeep list bookmarks "Reading List"
karakeep list create "New List"

# View bookmarks
karakeep show <bookmark_id>
karakeep recent [--limit 10]
karakeep tags                    # List all tags
```

#### 1.2 Batch Operations
```bash
# Bulk tagging
karakeep tag-add "python,tutorial" --search "python"

# Export bookmarks
karakeep export --format json > bookmarks.json
karakeep export --list "Work" --format csv

# Import from file
karakeep import urls.txt [--list "Imported"]

# Bulk operations
karakeep bulk-delete --tag "temp"
karakeep bulk-move --from "List A" --to "List B"
```

#### 1.3 Advanced Features
```bash
# Interactive mode
karakeep interactive     # Enter REPL-style interface

# Pipe support
cat urls.txt | karakeep add --stdin
karakeep search "api" | grep -i "python"

# Automation & scripting
karakeep sync-rss       # Sync RSS feeds
karakeep cleanup --older-than 30d --tag "temp"

# Statistics and analytics
karakeep stats                    # Overall statistics
karakeep stats --by-tag          # Tag usage
karakeep stats --by-date         # Timeline view
```

#### 1.4 Configuration
```bash
# Setup
karakeep config set-server https://my-karakeep.com
karakeep config set-token <jwt-token>

# Multiple profiles
karakeep config profile add work --server https://work.karakeep.com
karakeep config profile use work
karakeep config profile list
```

### Technical Architecture (CLI)

#### Technology Stack
- **Language:** Golang (Go)
  - **Why Go:** Single binary distribution, excellent performance, strong standard library, strict typing, excellent concurrency support.
  
- **HTTP Client:** `go-resty/resty` or standard `net/http`
- **CLI Framework:** `spf13/cobra` (standard for Go CLIs like kubectl, hugo)
- **Data Display:** 
  - `charmbracelet/lipgloss` for styling
  - `olekukonko/tablewriter` for tables
  - `charmbracelet/bubbletea` for interactive TUI elements
  
- **Configuration:** 
  - `spf13/viper` for config management
  - Store in `~/.karakeep/config.yaml`
  - Secure token storage using `99designs/keyring`

#### Project Structure
```
karakeep-cli/
├── cmd/
│   └── karakeep/      # Main entry point
│       └── main.go
├── internal/
│   ├── api/           # API client implementation
│   ├── auth/          # Authentication logic
│   ├── commands/      # CLI command handlers
│   ├── config/        # Configuration management
│   ├── models/        # Data structures
│   └── ui/            # Output formatting & TUI
├── pkg/               # Reusable library code (optional)
├── go.mod
├── go.sum
└── README.md
```

### CLI User Experience

#### Installation
```bash
# Install via Go
go install github.com/yourusername/karakeep-cli/cmd/karakeep@latest

# Or download pre-compiled binary
# (Future: brew install karakeep)

# First run setup wizard
karakeep setup
# Prompts for server URL and token
```

#### Output Formatting
- **Default:** Clean, readable terminal output
- **JSON mode:** `--json` flag for machine-readable output
- **Quiet mode:** `--quiet` for scripts (only errors)
- **Verbose mode:** `--verbose` for debugging

#### Error Handling
- Clear, actionable error messages
- Suggest fixes when possible
- Non-zero exit codes for scripting

---

## Phase 2: Web UI Application

### Differentiation Strategy

Rather than duplicating Karakeep's existing web interface, this application will focus on:

1. **Alternative Views & Visualizations**
   - Timeline view of bookmark history
   - Graph/network view of tag relationships
   - Gallery view for image bookmarks
   - Statistics dashboard with charts

2. **Specialized Workflows**
   - Reading queue management
   - Research project organization
   - Content curation tools
   - Bookmark deduplication

3. **Enhanced Search & Discovery**
   - Fuzzy search with typo tolerance
   - Advanced filters and saved searches
   - Related bookmarks suggestions
   - Bookmark recommendations based on patterns

4. **Collaboration Features** (if multi-user)
   - Shared lists with comments
   - Bookmark collections
   - Public bookmark pages

### Web UI Feature Set

#### 2.1 Dashboard
- Quick stats overview (total bookmarks, recent additions, top tags)
- Activity timeline
- Quick add bookmark form
- Saved search shortcuts

#### 2.2 Advanced Search Interface
- Visual query builder
- Search by multiple criteria simultaneously
- Saved searches with notifications
- Search result clustering by similarity

#### 2.3 Visualizations
- Tag cloud (interactive)
- Timeline view with filtering
- Network graph of tag relationships
- Heatmap of bookmark activity

#### 2.4 Bulk Management
- Multi-select with actions
- Merge duplicate bookmarks
- Clean up broken links
- Reorganize tags in bulk

#### 2.5 Analytics Dashboard
- Bookmark growth over time
- Most used tags
- Top domains saved
- Content type distribution
- Reading time estimates vs actual

### Technical Architecture (Web UI)

#### Frontend Stack
- **Framework:** React or Vue.js
  - **React:** Larger ecosystem, more job-relevant
  - **Vue.js:** Gentler learning curve, excellent docs
  
- **State Management:** 
  - React: Zustand or Context API
  - Vue: Pinia
  
- **UI Library:** 
  - Tailwind CSS for styling
  - shadcn/ui or Headless UI for components
  
- **Data Visualization:** 
  - Chart.js or Recharts for graphs
  - D3.js for advanced visualizations
  
- **HTTP Client:** axios or fetch API

#### Backend Options

**Option A: Serverless (Direct API)**
- Frontend directly calls Karakeep API
- No backend server needed
- Simpler deployment
- **Pros:** Simple, no server maintenance
- **Cons:** Limited to Karakeep API capabilities

**Option B: Proxy Server (Go Backend)**
- Lightweight Go server (Fiber or Gin)
- Proxies requests to Karakeep API
- Adds custom endpoints for specialized features
- **Pros:** Can implement custom logic, caching, rate limiting, single binary with embedded frontend
- **Cons:** Additional deployment complexity

**Recommendation:** Start with Option A, migrate to B if needed

#### Project Structure (Web)
```
karakeep-web/
├── src/
│   ├── components/        # React/Vue components
│   │   ├── Dashboard/
│   │   ├── Search/
│   │   ├── Visualizations/
│   │   └── Common/
│   ├── pages/            # Route pages
│   ├── services/         # API integration
│   │   └── api.js
│   ├── stores/           # State management
│   ├── utils/            # Helpers
│   ├── styles/           # CSS/styling
│   └── App.jsx
├── public/
├── package.json
└── README.md
```

#### Deployment Options
- **Static Hosting:** Netlify, Vercel, GitHub Pages
- **Self-Hosted:** Docker container alongside Karakeep
- **VPS:** Any server with Node.js or static file serving

---

## Development Roadmap

### Phase 1: CLI (Months 1-2)

#### Sprint 1: Foundation (Weeks 1-2)
- [ ] Set up Go project structure & modules
- [ ] Implement API client wrapper
- [ ] Basic authentication flow
- [ ] Configuration management (Viper)
- [ ] Basic commands: `add`, `search`, `list` (Cobra)

#### Sprint 2: Core Features (Weeks 3-4)
- [ ] Complete CRUD operations for bookmarks
- [ ] List management commands
- [ ] Tag operations
- [ ] Output formatting (Table, JSON)
- [ ] Error handling

#### Sprint 3: Advanced Features (Weeks 5-6)
- [ ] Batch operations
- [ ] Import/export functionality
- [ ] Search filters and advanced queries
- [ ] Interactive mode (Bubble Tea)
- [ ] Statistics commands

#### Sprint 4: Polish & Release (Weeks 7-8)
- [ ] Comprehensive testing
- [ ] Documentation
- [ ] Cross-platform compilation
- [ ] Release binary (GitHub Releases/Homebrew)

### Phase 2: Web UI (Months 3-5)

#### Sprint 1: Setup & Core (Weeks 9-11)
- [ ] Project scaffolding
- [ ] API integration layer
- [ ] Authentication flow
- [ ] Basic dashboard
- [ ] Bookmark list view

#### Sprint 2: Search & Filters (Weeks 12-14)
- [ ] Advanced search interface
- [ ] Filter builder
- [ ] Search results display
- [ ] Saved searches

#### Sprint 3: Visualizations (Weeks 15-17)
- [ ] Statistics dashboard
- [ ] Timeline view
- [ ] Tag visualization
- [ ] Analytics charts

#### Sprint 4: Advanced Features (Weeks 18-20)
- [ ] Bulk operations UI
- [ ] Bookmark deduplication
- [ ] Export functionality
- [ ] User preferences

#### Sprint 5: Polish & Deploy (Weeks 21-22)
- [ ] Responsive design
- [ ] Performance optimization
- [ ] User testing
- [ ] Documentation
- [ ] Deployment setup

---

## Technical Considerations

### Security
- **Token Storage:** Use secure keyring/keychain on CLI
- **HTTPS Only:** Enforce secure connections
- **Token Rotation:** Support updating tokens without reconfiguration
- **CORS:** Handle CORS for web UI if needed

### Performance
- **Concurrency:** Leverage Goroutines for concurrent API requests (e.g., batch operations)
- **Caching:** Implement smart caching for frequently accessed data
- **Pagination:** Support paginated results for large datasets
- **Lazy Loading:** Load data as needed in web UI

### Error Handling
- Network errors (connection issues, timeouts)
- Authentication errors (invalid/expired token)
- API errors (rate limiting, malformed requests)
- User input validation errors

### Testing Strategy
- **Unit Tests:** Go standard `testing` package, `testify` for assertions
- **Integration Tests:** Full command workflows
- **E2E Tests:** (Web) Full user journeys
- **Mock API:** Test without live Karakeep instance

---

## Monetization & Sustainability (Optional)

While this is an open-source tool, consider:

1. **Open Source (Free)**
   - Build reputation in self-hosted community
   - Accept donations via GitHub Sponsors
   - Offer premium support

2. **Freemium Model**
   - Basic features free
   - Advanced visualizations paid
   - Commercial license for businesses

3. **SaaS Option**
   - Host web UI for users who don't want to self-host
   - Subscription model
   - Connect to their Karakeep instance

---

## Risks & Mitigation

| Risk | Impact | Mitigation |
|------|--------|------------|
| Karakeep API changes | High | Version the client, maintain compatibility matrix |
| Limited API features | Medium | Implement workarounds, contribute to Karakeep core |
| User adoption | Medium | Focus on unique value proposition, marketing |
| Maintenance burden | Medium | Automate testing, CI/CD, clear documentation |
| Security vulnerabilities | High | Regular security audits, dependency updates |

---

## Success Metrics

### CLI Success Metrics
- Downloads/installations per month
- Active users (telemetry opt-in)
- GitHub stars and community engagement
- Issue resolution time

### Web UI Success Metrics
- Monthly active users
- Average session duration
- Feature usage statistics
- User feedback/satisfaction score

---

## Competitive Landscape

### Similar Tools
- **Karakeep Official CLI:** Basic functionality exists
- **Browser Extensions:** Quick saving, limited management
- **Raindrop.io:** Commercial, not self-hosted
- **Shiori:** Self-hosted, less features

### Our Differentiation
- **Deeper integration** with Karakeep's full API
- **Power user focus** with advanced CLI features
- **Unique visualizations** in web UI
- **Open source** and extensible
- **Automation-friendly** design

---

## Community & Contribution

### Open Source Strategy
- **License:** MIT or Apache 2.0
- **Repository:** GitHub with clear contribution guidelines
- **Documentation:** Comprehensive README, wiki, examples
- **Issues:** Bug reports, feature requests welcome
- **Discussions:** Community forum or Discord

### Integration with Karakeep Community
- Coordinate with Karakeep maintainer (Mohamed Bassem)
- Cross-promote in respective communities
- Consider upstreaming beneficial features
- Maintain compatibility with Karakeep releases

---

## Next Steps

### Immediate Actions (Week 1)
1. **Initialize Go module** and project repository
2. **Set up development environment** (Go 1.21+)
3. **Create Karakeep test instance** for development
4. **Generate API token** and test basic endpoints
5. **Prototype basic `add` and `search` commands** using Cobra

### Short Term (Weeks 2-4)
1. Implement core API wrapper
2. Build essential CLI commands
3. Set up testing framework
4. Create initial documentation

### Medium Term (Months 2-3)
1. Complete CLI feature set
2. Release CLI v1.0
3. Gather user feedback
4. Plan web UI architecture

---

## Resources & References

### Karakeep Resources
- **Documentation:** https://docs.karakeep.app
- **GitHub:** https://github.com/karakeep-app/karakeep
- **Demo Instance:** https://try.karakeep.app
- **API Docs:** https://docs.karakeep.app/api/karakeep-api

### Technical References
- **Cobra (Go CLI):** https://cobra.dev
- **Viper (Config):** https://github.com/spf13/viper
- **Resty (HTTP):** https://github.com/go-resty/resty
- **Lipgloss (Styles):** https://github.com/charmbracelet/lipgloss
- **Bubble Tea (TUI):** https://github.com/charmbracelet/bubbletea
- **React Documentation:** https://react.dev
- **Vue.js Guide:** https://vuejs.org/guide

### Community
- **Karakeep Discord:** https://discord.gg/NrgeYywsFh
- **r/selfhosted:** Reddit community for self-hosted apps

---

## Conclusion

This project offers an exciting opportunity to build a valuable tool for the Karakeep ecosystem. By starting with a focused CLI and evolving to a specialized web interface, we can deliver value incrementally while building a sustainable and maintainable application.

The key to success will be:
1. **Solving real problems** that Karakeep users face
2. **Maintaining quality** with good testing and documentation
3. **Engaging the community** for feedback and contributions
4. **Staying flexible** to adapt as Karakeep evolves

---

**Document Version:** 1.1
**Last Updated:** December 4, 2025
**Author:** Product Planning Team
**Status:** Draft - Ready for Review