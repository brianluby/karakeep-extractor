package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	gh "github.com/brianluby/karakeep-extractor/internal/adapter/github"
	"github.com/brianluby/karakeep-extractor/internal/adapter/http"
	"github.com/brianluby/karakeep-extractor/internal/adapter/karakeep"
	"github.com/brianluby/karakeep-extractor/internal/adapter/sqlite"
	"github.com/brianluby/karakeep-extractor/internal/adapter/trillium"
	"github.com/brianluby/karakeep-extractor/internal/config"
	"github.com/brianluby/karakeep-extractor/internal/core/domain"
	"github.com/brianluby/karakeep-extractor/internal/core/service"
	"github.com/brianluby/karakeep-extractor/internal/ui"
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

	rankCmd := flag.NewFlagSet("rank", flag.ExitOnError)
	rankLimit := rankCmd.Int("limit", 20, "Number of repositories to display")
	rankSort := rankCmd.String("sort", "stars", "Metric to sort by (stars, forks, updated)")
	rankFormat := rankCmd.String("format", "table", "Output format (table, json, csv)")
	rankSinkURL := rankCmd.String("sink-url", "", "URL to POST ranked results to")
	var rankSinkHeaders arrayFlags
	rankCmd.Var(&rankSinkHeaders, "sink-header", "Header to send with sink request (Key: Value)")
	rankSinkTrillium := rankCmd.Bool("sink-trillium", false, "Send ranked results to Trillium Notes")
	rankTag := rankCmd.String("tag", "", "Filter repositories by tag (title/description)")

	// Global flags logic is complex with subcommands if mixed. 
	// We'll assume extract is default if no subcommand, or explicit 'extract' command.
	// For now, let's support "extract" and "enrich" explicitly.
	
	command := os.Args[1]

	switch command {
	case "extract":
		runExtract()
	case "enrich":
		// Parse flags for enrich
		enrichCmd.Parse(os.Args[2:])
		runEnrich(*enrichLimit, *enrichForce, *enrichToken)
	case "rank":
		rankCmd.Parse(os.Args[2:])
		runRank(*rankLimit, *rankSort, *rankFormat, *rankSinkURL, rankSinkHeaders, *rankSinkTrillium, *rankTag)
	case "setup":
		// Fallback to extract for backward compatibility or print usage?
		// Plan implied "karakeep enrich" as a command.
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Karakeep Extractor - Intelligence for your bookmarks")
	fmt.Println("")
	fmt.Println("Usage: karakeep <command> [flags]")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  setup      Run the interactive configuration wizard to set API tokens and URLs.")
	fmt.Println("  extract    Fetch bookmarks from Karakeep and save GitHub links to the local database.")
	fmt.Println("  enrich     Fetch metadata (stars, forks, etc.) from GitHub for extracted repositories.")
	fmt.Println("  rank       Display, filter, and export a ranked list of repositories.")
	fmt.Println("")
	fmt.Println("Run 'karakeep <command> --help' for command-specific flags.")
}

func runExtract() {
	// Create a dedicated FlagSet for extract to parse args after the subcommand
	extractCmd := flag.NewFlagSet("extract", flag.ExitOnError)
	var (
		karakeepURL   = extractCmd.String("url", "", "Karakeep Base URL")
		karakeepToken = extractCmd.String("token", "", "Karakeep API Token")
		dbPath        = extractCmd.String("db", "./karakeep.db", "Path to SQLite database")
	)

	// Parse arguments starting from os.Args[2]
	extractCmd.Parse(os.Args[2:])

	// Load Env vars as fallback
	if *karakeepURL == "" {
		*karakeepURL = os.Getenv("KARAKEEP_URL")
	}
	if *karakeepToken == "" {
		*karakeepToken = os.Getenv("KARAKEEP_TOKEN")
	}
	if *dbPath == "./karakeep.db" && os.Getenv("KARAKEEP_DB") != "" {
		*dbPath = os.Getenv("KARAKEEP_DB")
	}

	if *karakeepURL == "" || *karakeepToken == "" {
		fmt.Println("Error: Karakeep URL and Token are required.")
		os.Exit(1)
	}

	domainCfg := &domain.KarakeepConfig{
		BaseURL:  *karakeepURL,
		APIToken: *karakeepToken,
	}
	client := karakeep.NewClient(domainCfg)

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
	if err := svc.Extract(context.Background()); err != nil {
		log.Fatalf("Extraction failed: %v", err)
	}
}

func runEnrich(limit int, force bool, tokenOverride string) {
	// Load base config for DB path
	// We can't use config.Load() easily because it calls flag.Parse() which acts on global args.
	// We'll just grab DB from Env or Default.
	dbPath := os.Getenv("KARAKEEP_DB")
	if dbPath == "" {
		dbPath = "./karakeep.db"
	}
	
	ghToken := os.Getenv("GITHUB_TOKEN")
	if tokenOverride != "" {
		ghToken = tokenOverride
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

	fmt.Printf("Starting enrichment (Limit: %d, Force: %t)...\n", limit, force)
	success, failed, err := enricher.EnrichBatch(context.Background(), limit, force, 5) // 5 workers
	
	fmt.Printf("Enrichment complete.\nUpdated: %d\nErrors: %d\n", success, failed)
	
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runRank(limit int, sort string, format string, sinkURL string, sinkHeaders []string, sinkTrillium bool, tag string) {
	// Load Config
	loader := config.NewConfigLoader()
	cfg, err := loader.LoadConfig(nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to load config file: %v\n", err)
	}

	db, err := sql.Open("sqlite3", cfg.DBPath)
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
			fmt.Fprintf(os.Stderr, "Error: Trillium URL and Token required. Run 'karakeep setup'.\n")
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
