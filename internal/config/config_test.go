package config

import (
	"flag"
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Save original args and env
	origArgs := os.Args
	origEnvURL := os.Getenv("KARAKEEP_URL")
	origEnvToken := os.Getenv("KARAKEEP_TOKEN")
	defer func() {
		os.Args = origArgs
		os.Setenv("KARAKEEP_URL", origEnvURL)
		os.Setenv("KARAKEEP_TOKEN", origEnvToken)
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError) // Reset flags
	}()

	// Test case 1: Env vars only
	os.Setenv("KARAKEEP_URL", "http://env.com")
	os.Setenv("KARAKEEP_TOKEN", "env-token")
	os.Args = []string{"cmd"}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	cfg := Load()
	if cfg.KarakeepURL != "http://env.com" {
		t.Errorf("Expected URL http://env.com, got %s", cfg.KarakeepURL)
	}
	if cfg.KarakeepToken != "env-token" {
		t.Errorf("Expected Token env-token, got %s", cfg.KarakeepToken)
	}

	// Test case 2: Flags override Env
	os.Args = []string{"cmd", "--url", "http://flag.com", "--token", "flag-token"}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	cfg = Load()
	if cfg.KarakeepURL != "http://flag.com" {
		t.Errorf("Expected URL http://flag.com, got %s", cfg.KarakeepURL)
	}
	if cfg.KarakeepToken != "flag-token" {
		t.Errorf("Expected Token flag-token, got %s", cfg.KarakeepToken)
	}
}
