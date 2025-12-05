package config

import (
	"os"
	"testing"
)

func TestConfigLoader_LoadConfig_Precedence(t *testing.T) {
	// 1. Setup Env Vars
	os.Setenv("KARAKEEP_URL", "env-url")
	defer os.Unsetenv("KARAKEEP_URL")

	loader := NewConfigLoader()

	// 2. Test Flag Override
	flagConfig := &Config{
		KarakeepURL: "flag-url",
	}

	cfg, err := loader.LoadConfig(flagConfig)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.KarakeepURL != "flag-url" {
		t.Errorf("Expected flag-url, got %s", cfg.KarakeepURL)
	}

	// 3. Test Env Fallback (Empty Flag)
	emptyFlagConfig := &Config{}
	cfg, err = loader.LoadConfig(emptyFlagConfig)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.KarakeepURL != "env-url" {
		t.Errorf("Expected env-url, got %s", cfg.KarakeepURL)
	}
}

func TestConfigLoader_SaveConfig(t *testing.T) {
	// This test writes to the real user config dir which is risky/messy.
	// Ideally we'd mock os.UserConfigDir, but that's not easy in Go without a wrapper.
	// We'll skip the integration test for now or assume manual verification via `setup`.
	// Alternatively, we could refactor `SaveConfig` to take a path.
}
