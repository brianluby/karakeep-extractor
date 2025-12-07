package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	gh "github.com/brianluby/karakeep-extractor/internal/adapter/github"
	"github.com/brianluby/karakeep-extractor/internal/adapter/http"
	"github.com/brianluby/karakeep-extractor/internal/adapter/karakeep"
	"github.com/brianluby/karakeep-extractor/internal/adapter/llm"
	rep "github.com/brianluby/karakeep-extractor/internal/adapter/reporter"
	"github.com/brianluby/karakeep-extractor/internal/adapter/sqlite"
	"github.com/brianluby/karakeep-extractor/internal/adapter/trillium"
	"github.com/brianluby/karakeep-extractor/internal/config"
	"github.com/brianluby/karakeep-extractor/internal/core/domain"
	"github.com/brianluby/karakeep-extractor/internal/core/service"
	"github.com/brianluby/karakeep-extractor/internal/core/service/analysis"
	"github.com/brianluby/karakeep-extractor/internal/ui"
	"github.com/brianluby/karakeep-extractor/internal/ui/tui"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	// Subcommand handling
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Define subcommands
	enrichCmd := flag.NewFlagSet("enrich", flag.ExitOnError)
	enrichLimit := enrichCmd.Int("limit", 50, "Maximum number of repositories to process")
	enrichForce := enrichCmd.Bool("force", false, "Force re-enrichment of processed repositories")
	enrichToken := enrichCmd.String("token", "", "GitHub Personal Access Token (overrides env var)")
	enrichDB := enrichCmd.String("db", "", "Path to SQLite database")
	enrichTui := enrichCmd.Bool("tui", false, "Enable TUI mode")

	rankCmd := flag.NewFlagSet("rank", flag.ExitOnError)
	rankLimit := rankCmd.Int("limit", 20, "Number of repositories to display")
	rankSort := rankCmd.String("sort", "stars", "Metric to sort by (stars, forks, updated)")
	rankFormat := rankCmd.String("format", "table", "Output format (table, json, csv)")
	rankSinkURL := rankCmd.String("sink-url", "", "URL to POST ranked results to")
	var rankSinkHeaders arrayFlags
	rankCmd.Var(&rankSinkHeaders, "sink-header", "Header to send with sink request (Key: Value)")
	rankSinkTrillium := rankCmd.Bool("sink-trillium", false, "Send ranked results to Trillium Notes")
	rankTag := rankCmd.String("tag", "", "Filter repositories by tag (title/description)")
	rankDB := rankCmd.String("db", "", "Path to SQLite database")

	analyzeCmd := flag.NewFlagSet("analyze", flag.ExitOnError)
	analyzeLang := analyzeCmd.String("lang", "", "Filter by language")
	analyzeLimit := analyzeCmd.Int("limit", 50, "Limit number of repositories")
	analyzeTag := analyzeCmd.String("tag", "", "Filter by tag")
	analyzeDB := analyzeCmd.String("db", "", "Path to SQLite database")
	analyzeMinStars := analyzeCmd.Int("min-stars", 0, "Minimum number of stars")
	analyzeMaxStars := analyzeCmd.Int("max-stars", 0, "Maximum number of stars (0 for no limit)")

	// Global flags logic is complex with subcommands if mixed. 
	// We'll assume extract is default if no subcommand, or explicit 'extract' command.
	// For now, let's support "extract" and "enrich" explicitly.
	
	command := os.Args[1]

	switch command {
	case "--help", "-h":
		printUsage()
		os.Exit(0)
	case "extract":
		runExtract()
	case "enrich":
		// Parse flags for enrich
		enrichCmd.Parse(os.Args[2:])
		runEnrich(*enrichLimit, *enrichForce, *enrichToken, *enrichDB, *enrichTui)
	case "rank":
		rankCmd.Parse(os.Args[2:])
		runRank(*rankLimit, *rankSort, *rankFormat, *rankSinkURL, rankSinkHeaders, *rankSinkTrillium, *rankTag, *rankDB)
	case "setup":
		runSetup()
	case "config":
		if len(os.Args) < 3 || os.Args[2] != "llm" {
			fmt.Println("Usage: karakeep config llm")
			os.Exit(1)
		}
		runConfigLLM()
	case "analyze":
		analyzeCmd.Parse(os.Args[2:])
		if analyzeCmd.NArg() < 1 {
			fmt.Println("Usage: karakeep analyze [flags] \"query\"")
			os.Exit(1)
		}
		query := analyzeCmd.Arg(0)
		runAnalyze(*analyzeLang, *analyzeLimit, *analyzeTag, *analyzeDB, *analyzeMinStars, *analyzeMaxStars, query)
	}
}

func printUsage() {
	fmt.Println("Karakeep Extractor - Intelligence for your bookmarks")
	fmt.Println("")
	fmt.Println("Usage: karakeep-extractor <command> [flags]")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  setup      Run the interactive configuration wizard to set API tokens and URLs.")
	fmt.Println("  config     Manage configuration (e.g., 'config llm').")
	fmt.Println("  extract    Fetch bookmarks from Karakeep and save GitHub links to the local database.")
	fmt.Println("  enrich     Fetch metadata (stars, forks, etc.) from GitHub for extracted repositories.")
	fmt.Println("  rank       Display, filter, and export a ranked list of repositories.")
	fmt.Println("  analyze    Analyze repositories using an LLM.")
	fmt.Println("")
	fmt.Println("Run 'karakeep-extractor <command> --help' for command-specific flags.")
}

// expandPath expands the tilde (~) in the path to the user's home directory.
func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err == nil {
			return filepath.Join(home, path[2:])
		}
	}
	return path
}

func runConfigLLM() {
	prompt := ui.NewPrompt(os.Stdin, os.Stdout)
	loader := config.NewConfigLoader()
	currentCfg, _ := loader.LoadConfig(nil)

	fmt.Println("Karakeep LLM Configuration")
	fmt.Println("--------------------------")

	// Defaults
	defaultProvider := "openai"
	if currentCfg.LLM.Provider != "" {
		defaultProvider = currentCfg.LLM.Provider
	}
	defaultBaseURL := "https://api.openai.com/v1"
	if currentCfg.LLM.BaseURL != "" {
		defaultBaseURL = currentCfg.LLM.BaseURL
	}
	defaultModel := "gpt-4o"
	if currentCfg.LLM.Model != "" {
		defaultModel = currentCfg.LLM.Model
	}

	// Prompts
	provider, err := prompt.Ask("Provider (openai, anthropic, local)", defaultProvider)
	if err != nil {
		log.Fatalf("Error reading input: %v", err)
	}

	baseURL, err := prompt.Ask("Base URL", defaultBaseURL)
	if err != nil {
		log.Fatalf("Error reading input: %v", err)
	}

	apiKey, err := prompt.AskSecret("API Key")
	if err != nil {
		log.Fatalf("Error reading input: %v", err)
	}
	if apiKey == "" {
		apiKey = currentCfg.LLM.APIKey
	}

	model, err := prompt.Ask("Model Name", defaultModel)
	if err != nil {
		log.Fatalf("Error reading input: %v", err)
	}

	// Update Config
	currentCfg.LLM.Provider = provider
	currentCfg.LLM.BaseURL = baseURL
	currentCfg.LLM.APIKey = apiKey
	currentCfg.LLM.Model = model

	// Save
	if err := loader.SaveConfig(currentCfg); err != nil {
		log.Fatalf("Failed to save config: %v", err)
	}
	
	path, _ := config.GetConfigPath()
	fmt.Printf("\nLLM configuration saved to %s\n", path)
}

func runAnalyze(lang string, limit int, tag string, dbFlag string, minStars int, maxStars int, query string) {
	// 1. Config
	loader := config.NewConfigLoader()
	cfg, err := loader.LoadConfig(nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to load config: %v\n", err)
	}

	// Basic Validation
	if cfg.LLM.BaseURL == "" {
		fmt.Println("Error: LLM not configured. Run 'karakeep config llm'.")
		os.Exit(1)
	}

	// 2. DB
	dbPath := dbFlag
	if dbPath == "" {
		dbPath = os.Getenv("KARAKEEP_DB")
	}
	if dbPath == "" && cfg != nil {
		dbPath = cfg.DBPath
	}
	if dbPath == "" {
		dbPath = "./karakeep.db"
	}
	dbPath = expandPath(dbPath)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}
	defer db.Close()
	repo := sqlite.NewSQLiteRepository(db)
	
	// 3. Service
	llmClient := llm.NewClient(cfg.LLM)
	svc := analysis.NewService(repo, llmClient)

	fmt.Println("Analyzing repositories...")
	answer, err := svc.Analyze(context.Background(), query, limit, lang, tag, minStars, maxStars)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error during analysis: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n--- Analysis Result ---")
	fmt.Println(answer)
}

func runExtract() {
	// Load Config first to get defaults from file
	loader := config.NewConfigLoader()
	cfg, err := loader.LoadConfig(nil)
	if err != nil {
		// It's okay if config file doesn't exist, we'll fallback to flags/env
		// but if it exists and is malformed, maybe warn?
	}

	// Create a dedicated FlagSet for extract to parse args after the subcommand
	extractCmd := flag.NewFlagSet("extract", flag.ExitOnError)
	var (
		karakeepURL   = extractCmd.String("url", "", "Karakeep Base URL")
		karakeepToken = extractCmd.String("token", "", "Karakeep API Token")
		dbPath        = extractCmd.String("db", "", "Path to SQLite database")
		tuiMode       = extractCmd.Bool("tui", false, "Enable TUI mode")
	)

	// Parse arguments starting from os.Args[2]
	extractCmd.Parse(os.Args[2:])

	// Precedence: Flag > Env > Config File > Default

	// 1. URL
	if *karakeepURL == "" {
		*karakeepURL = os.Getenv("KARAKEEP_URL")
	}
	if *karakeepURL == "" && cfg != nil {
		*karakeepURL = cfg.KarakeepURL
	}

	// 2. Token
	if *karakeepToken == "" {
		*karakeepToken = os.Getenv("KARAKEEP_TOKEN")
	}
	if *karakeepToken == "" && cfg != nil {
		*karakeepToken = cfg.KarakeepToken
	}

	// 3. DB Path
	if *dbPath == "" {
		*dbPath = os.Getenv("KARAKEEP_DB")
	}
	if *dbPath == "" && cfg != nil {
		*dbPath = cfg.DBPath
	}
	if *dbPath == "" {
		*dbPath = "./karakeep.db" // Final default
	}
	
	*dbPath = expandPath(*dbPath)

	if *karakeepURL == "" || *karakeepToken == "" {
		fmt.Println("Error: Karakeep URL and Token are required.")
		fmt.Println("Run 'karakeep-extractor setup' or provide via flags/env.")
		os.Exit(1)
	}

	domainCfg := &domain.KarakeepConfig{
		BaseURL:  *karakeepURL,
		APIToken: *karakeepToken,
	}
	client := karakeep.NewClient(domainCfg)

	// Ensure DB directory exists
	dbDir := filepath.Dir(*dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		log.Fatalf("Failed to create DB directory: %v", err)
	}

	db, err := sql.Open("sqlite3", *dbPath)
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}
	defer db.Close()
	repo := sqlite.NewSQLiteRepository(db)
	if err := repo.InitSchema(context.Background()); err != nil {
		log.Fatalf("Schema init failed: %v", err)
	}

	svc := service.NewExtractor(client, repo)

	// Select Reporter
	var reporter domain.ProgressReporter
	if *tuiMode {
		task := func(r domain.ProgressReporter) error {
			return svc.Extract(context.Background(), r)
		}

		// Run TUI
		if err := tui.Run(context.Background(), "extract", task); err != nil {
			fmt.Fprintf(os.Stderr, "TUI Error: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	} else {
		reporter = rep.NewTextReporter()
		if err := svc.Extract(context.Background(), reporter); err != nil {
			log.Fatalf("Extraction failed: %v", err)
		}
	}
}

func runEnrich(limit int, force bool, tokenOverride string, dbFlag string, tuiMode bool) {
	// Load Config
	loader := config.NewConfigLoader()
	cfg, err := loader.LoadConfig(nil)
	// Ignore err, just try to get values

	// DB Path Precedence: Flag > Env > Config > Default
	dbPath := dbFlag
	if dbPath == "" {
		dbPath = os.Getenv("KARAKEEP_DB")
	}
	if dbPath == "" && cfg != nil {
		dbPath = cfg.DBPath
	}
	if dbPath == "" {
		dbPath = "./karakeep.db"
	}
	
	dbPath = expandPath(dbPath)
	
	// GitHub Token Precedence: Flag > Env > Config
	ghToken := tokenOverride
	if ghToken == "" {
		ghToken = os.Getenv("GITHUB_TOKEN")
	}
	if ghToken == "" && cfg != nil {
		ghToken = cfg.GitHubToken
	}

	// Ensure DB directory exists
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		log.Fatalf("Failed to create DB directory: %v", err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}
	defer db.Close()
	repo := sqlite.NewSQLiteRepository(db)
	// Ensure schema is up to date (migrations)
	if err := repo.InitSchema(context.Background()); err != nil {
		log.Fatalf("Schema init failed: %v", err)
	}

	ghClient := gh.NewClient(ghToken)
	enricher := service.NewEnricher(repo, ghClient)

	// Select Reporter
	var reporter domain.ProgressReporter
	if tuiMode {
		task := func(r domain.ProgressReporter) error {
			// fmt.Printf("Starting enrichment (Limit: %d, Force: %t)...\n", limit, force) // Handled by Reporter
			_, _, err := enricher.EnrichBatch(context.Background(), limit, force, 5, r) // 5 workers
			return err
		}

		// Run TUI
		if err := tui.Run(context.Background(), "enrich", task); err != nil {
			fmt.Fprintf(os.Stderr, "TUI Error: %v\n", err)
			os.Exit(1)
		}
		// Exit successfully if TUI loop finishes normally
		os.Exit(0)
	} else {
		reporter = rep.NewTextReporter()
		success, failed, err := enricher.EnrichBatch(context.Background(), limit, force, 5, reporter) // 5 workers
		if err != nil {
			os.Exit(1)
		}
		// If using text reporter, we might want to log summary if not already done by Finish()
		// TextReporter implementation does log "Finished: ...".
		_ = success
		_ = failed
	}
}

func runRank(limit int, sort string, format string, sinkURL string, sinkHeaders []string, sinkTrillium bool, tag string, dbFlag string) {
	// Load Config
	loader := config.NewConfigLoader()
	cfg, err := loader.LoadConfig(nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to load config file: %v\n", err)
	}
	
	// DB Path Precedence: Flag > Env > Config > Default
	dbPath := dbFlag
	if dbPath == "" {
		dbPath = os.Getenv("KARAKEEP_DB")
	}
	if dbPath == "" && cfg != nil {
		dbPath = cfg.DBPath
	}
	if dbPath == "" {
		dbPath = "./karakeep.db"
	}
	
	dbPath = expandPath(dbPath)

	// Ensure DB directory exists (safe measure)
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		log.Fatalf("Failed to create DB directory: %v", err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}
	defer db.Close()
	repo := sqlite.NewSQLiteRepository(db)
	// Schema init not strictly required if just reading, but good practice
	if err := repo.InitSchema(context.Background()); err != nil {
		log.Fatalf("Schema init failed: %v", err)
	}

	var exporter domain.Exporter
	exporter, err = ui.GetExporter(format)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	var sink domain.Sink
	if sinkTrillium {
		if cfg.TrilliumURL == "" || cfg.TrilliumToken == "" {
			fmt.Fprintf(os.Stderr, "Error: Trillium URL and Token required. Run 'karakeep-extractor setup'.\n")
			os.Exit(1)
		}
		client := trillium.NewClient(cfg.TrilliumURL, cfg.TrilliumToken)
		sink = trillium.NewSink(client)
	} else if sinkURL != "" {
		sink = http.NewHTTPSink(sinkURL, sinkHeaders)
	}

	ranker := service.NewRanker(repo, exporter, sink)
	if err := ranker.Rank(context.Background(), limit, sort, tag, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runSetup() {
	prompt := ui.NewPrompt(os.Stdin, os.Stdout)

	fmt.Println("Karakeep Extractor Setup")
	fmt.Println("------------------------")

	// 1. Load existing or default
	loader := config.NewConfigLoader()
	
	// We load config (ignoring flags since setup is clean interaction usually, but we could respect them as defaults)
	// For now, just load file/env to populate defaults
	currentCfg, _ := loader.LoadConfig(nil)

	// 2. Prompt User
	url, err := prompt.Ask("Enter Karakeep URL", currentCfg.KarakeepURL)
	if err != nil {
		log.Fatalf("Error reading input: %v", err)
	}

	token, err := prompt.AskSecret("Enter Karakeep API Token")
	if err != nil {
		log.Fatalf("Error reading input: %v", err)
	}
	// If user just presses enter on secret, it returns empty string.
	// We assume they want to keep existing token if they provide empty.
	if token == "" {
		token = currentCfg.KarakeepToken
	}

	ghToken, err := prompt.AskSecret("Enter GitHub Personal Access Token (optional)")
	if err != nil {
		log.Fatalf("Error reading input: %v", err)
	}
	if ghToken == "" {
		ghToken = currentCfg.GitHubToken
	}

	dbPath, err := prompt.Ask("Enter SQLite Database Path", currentCfg.DBPath)
	if err != nil {
		log.Fatalf("Error reading input: %v", err)
	}

	// Integrations
	configureTrillium, err := prompt.AskConfirm("Configure Trillium Integration?")
	if err != nil {
		log.Fatalf("Error reading input: %v", err)
	}

	var trilliumURL, trilliumToken string
	if configureTrillium {
		trilliumURL, err = prompt.Ask("Enter Trillium Instance URL", currentCfg.TrilliumURL)
		if err != nil {
			log.Fatalf("Error reading input: %v", err)
		}
		
		trilliumToken, err = prompt.AskSecret("Enter Trillium ETAPI Token")
		if err != nil {
			log.Fatalf("Error reading input: %v", err)
		}
		if trilliumToken == "" {
			trilliumToken = currentCfg.TrilliumToken
		}
	} else {
		// Keep existing if not re-configuring? Or clear?
		// Typically if user says "No" to configuring, we might skip changing it.
		// But if they want to disable it?
		// For MVP setup, let's assume we keep existing values if they skip the section, 
		// or we could just use the current values as defaults and they can clear them?
		// The prompt logic above uses currentCfg as defaults.
		// If they skip the section, we just use current values?
		trilliumURL = currentCfg.TrilliumURL
		trilliumToken = currentCfg.TrilliumToken
	}

	// 3. Confirm Overwrite if file exists
	path, _ := config.GetConfigPath()
	if _, err := os.Stat(path); err == nil {
		confirm, _ := prompt.Ask(fmt.Sprintf("Overwrite existing config at %s? (y/N)", path), "N")
		if strings.ToLower(confirm) != "y" {
			fmt.Println("Aborted.")
			os.Exit(0)
		}
	}

	// 4. Save
	newCfg := &config.Config{
		KarakeepURL:   url,
		KarakeepToken: token,
		GitHubToken:   ghToken,
		DBPath:        dbPath,
		TrilliumURL:   trilliumURL,
		TrilliumToken: trilliumToken,
	}

	if err := loader.SaveConfig(newCfg); err != nil {
		log.Fatalf("Failed to save config: %v", err)
	}

	fmt.Printf("\nConfiguration saved to %s\n", path)
	fmt.Println("Permissions set to 0600.")
}
