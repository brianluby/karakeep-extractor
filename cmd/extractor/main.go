package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/brianluby/karakeep-extractor/internal/adapter/github"
	"github.com/brianluby/karakeep-extractor/internal/adapter/karakeep"
	"github.com/brianluby/karakeep-extractor/internal/adapter/sqlite"
	"github.com/brianluby/karakeep-extractor/internal/core/domain"
	"github.com/brianluby/karakeep-extractor/internal/core/service"
)

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
		runRank(*rankLimit, *rankSort)
	default:
		// Fallback to extract for backward compatibility or print usage?
		// Plan implied "karakeep enrich" as a command.
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: karakeep <command> [flags]")
	fmt.Println("Commands:")
	fmt.Println("  extract    Run the extraction process")
	fmt.Println("  enrich     Enrich extracted repositories with GitHub metadata")
	fmt.Println("  rank       Display ranked list of repositories")
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

	ghClient := github.NewClient(ghToken)
	enricher := service.NewEnricher(repo, ghClient)

	fmt.Printf("Starting enrichment (Limit: %d, Force: %t)...\n", limit, force)
	success, failed, err := enricher.EnrichBatch(context.Background(), limit, force, 5) // 5 workers
	
	fmt.Printf("Enrichment complete.\nUpdated: %d\nErrors: %d\n", success, failed)
	
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runRank(limit int, sort string) {
	dbPath := os.Getenv("KARAKEEP_DB")
	if dbPath == "" {
		dbPath = "./karakeep.db"
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

	ranker := service.NewRanker(repo)
	if err := ranker.Rank(context.Background(), limit, sort, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
