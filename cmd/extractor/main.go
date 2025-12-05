package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/brianluby/karakeep-extractor/internal/adapter/karakeep"
	"github.com/brianluby/karakeep-extractor/internal/adapter/sqlite"
	"github.com/brianluby/karakeep-extractor/internal/config"
	"github.com/brianluby/karakeep-extractor/internal/core/domain"
	"github.com/brianluby/karakeep-extractor/internal/core/service"
)

func main() {
	cfg := config.Load()

	if cfg.KarakeepURL == "" || cfg.KarakeepToken == "" {
		fmt.Println("Error: Karakeep URL and Token are required.")
		fmt.Println("Please provide them via --url and --token flags or KARAKEEP_URL and KARAKEEP_TOKEN environment variables.")
		os.Exit(1)
	}

	// Initialize Karakeep Client
	domainCfg := &domain.KarakeepConfig{
		BaseURL:  cfg.KarakeepURL,
		APIToken: cfg.KarakeepToken,
	}
	karakeepClient := karakeep.NewClient(domainCfg)

	// Create a context with a timeout for initial operations
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Printf("Pinging Karakeep instance at %s...", cfg.KarakeepURL)
	if err := karakeepClient.Ping(ctx); err != nil {
		log.Fatalf("Failed to connect to Karakeep: %v", err)
	}
	log.Println("Successfully connected to Karakeep.")

	// Initialize SQLite Database and Repository
	db, err := sql.Open("sqlite3", cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to open SQLite database at %s: %v", cfg.DBPath, err)
	}
	defer db.Close()

	repo := sqlite.NewSQLiteRepository(db)
	if err := repo.InitSchema(ctx); err != nil {
		log.Fatalf("Failed to initialize database schema: %v", err)
	}
	log.Printf("SQLite database initialized at %s", cfg.DBPath)

	// Initialize Extractor Service
	extractorService := service.NewExtractor(karakeepClient, repo)

	log.Println("Starting extraction process...")
	if err := extractorService.Extract(ctx); err != nil {
		log.Fatalf("Extraction failed: %v", err)
	}
	log.Println("Extraction process completed successfully.")

	// TODO: Add summary statistics (e.g., total repos extracted)
	os.Exit(0)
}
