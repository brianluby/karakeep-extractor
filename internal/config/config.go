package config

import (
	"flag"
	"os"
)

type Config struct {
	KarakeepURL   string
	KarakeepToken string
	DBPath        string
	GitHubToken   string
}

func Load() *Config {
	cfg := &Config{}

	// Define flags
	flag.StringVar(&cfg.KarakeepURL, "url", "", "Karakeep Base URL")
	flag.StringVar(&cfg.KarakeepToken, "token", "", "Karakeep API Token")
	flag.StringVar(&cfg.DBPath, "db", "./karakeep.db", "Path to SQLite database")
	flag.StringVar(&cfg.GitHubToken, "github-token", "", "GitHub Personal Access Token")

	// Parse flags
	flag.Parse()

	// Apply environment variables if flags are not set
	if cfg.KarakeepURL == "" {
		cfg.KarakeepURL = os.Getenv("KARAKEEP_URL")
	}
	if cfg.KarakeepToken == "" {
		cfg.KarakeepToken = os.Getenv("KARAKEEP_TOKEN")
	}
	if cfg.DBPath == "./karakeep.db" && os.Getenv("KARAKEEP_DB") != "" {
		cfg.DBPath = os.Getenv("KARAKEEP_DB")
	}
	if cfg.GitHubToken == "" {
		cfg.GitHubToken = os.Getenv("GITHUB_TOKEN")
	}

	return cfg
}
