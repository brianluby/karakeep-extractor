package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const ConfigDirName = "karakeep"
const ConfigFileName = "config.yaml"

// ConfigLoader handles loading and saving configuration with precedence.
type ConfigLoader struct{}

func NewConfigLoader() *ConfigLoader {
	return &ConfigLoader{}
}

// LoadConfig loads configuration with precedence: Flag > Env > File.
// flagConfig contains values parsed from CLI flags (if any, usually minimal in this new flow or passed explicitly).
// However, standard flag.Parse() populates variables.
// We need a way to merge.
// Strategy:
// 1. Load from File (if exists)
// 2. Override with Env
// 3. Override with Flags (passed in as argument)
func (l *ConfigLoader) LoadConfig(flagConfig *Config) (*Config, error) {
	// 1. Start with defaults or empty
	finalConfig := &Config{
		DBPath: "./karakeep.db",
	}

	// 2. Load from File
	configPath, err := GetConfigPath()
	if err == nil {
		if _, err := os.Stat(configPath); err == nil {
			fileData, err := os.ReadFile(configPath)
			if err != nil {
				return nil, fmt.Errorf("failed to read config file: %w", err)
			}
			var fileConfig Config
			if err := yaml.Unmarshal(fileData, &fileConfig); err != nil {
				return nil, fmt.Errorf("failed to parse config file: %w", err)
			}
			// Merge file config
			if fileConfig.KarakeepURL != "" {
				finalConfig.KarakeepURL = fileConfig.KarakeepURL
			}
			if fileConfig.KarakeepToken != "" {
				finalConfig.KarakeepToken = fileConfig.KarakeepToken
			}
			if fileConfig.GitHubToken != "" {
				finalConfig.GitHubToken = fileConfig.GitHubToken
			}
			if fileConfig.TrilliumURL != "" {
				finalConfig.TrilliumURL = fileConfig.TrilliumURL
			}
			if fileConfig.TrilliumToken != "" {
				finalConfig.TrilliumToken = fileConfig.TrilliumToken
			}
			if fileConfig.DBPath != "" {
				finalConfig.DBPath = fileConfig.DBPath
			}
		}
	}

	// 3. Override with Env
	if val := os.Getenv("KARAKEEP_URL"); val != "" {
		finalConfig.KarakeepURL = val
	}
	if val := os.Getenv("KARAKEEP_TOKEN"); val != "" {
		finalConfig.KarakeepToken = val
	}
	if val := os.Getenv("GITHUB_TOKEN"); val != "" {
		finalConfig.GitHubToken = val
	}
	if val := os.Getenv("TRILLIUM_URL"); val != "" {
		finalConfig.TrilliumURL = val
	}
	if val := os.Getenv("TRILLIUM_TOKEN"); val != "" {
		finalConfig.TrilliumToken = val
	}
	if val := os.Getenv("KARAKEEP_DB"); val != "" {
		finalConfig.DBPath = val
	}

	// 4. Override with Flags (if provided)
	if flagConfig != nil {
		if flagConfig.KarakeepURL != "" {
			finalConfig.KarakeepURL = flagConfig.KarakeepURL
		}
		if flagConfig.KarakeepToken != "" {
			finalConfig.KarakeepToken = flagConfig.KarakeepToken
		}
		if flagConfig.GitHubToken != "" {
			finalConfig.GitHubToken = flagConfig.GitHubToken
		}
		if flagConfig.TrilliumURL != "" {
			finalConfig.TrilliumURL = flagConfig.TrilliumURL
		}
		if flagConfig.TrilliumToken != "" {
			finalConfig.TrilliumToken = flagConfig.TrilliumToken
		}
		if flagConfig.DBPath != "" && flagConfig.DBPath != "./karakeep.db" { // Basic check to see if flag was actually set differently from default
			// Note: This is imperfect if the flag default matches our internal default.
			// A robust solution would use a library like `viper` or explicit "isSet" flags.
			// For now, we assume if it's not empty/default, it overrides.
			finalConfig.DBPath = flagConfig.DBPath
		}
	}

	return finalConfig, nil
}

// SaveConfig writes the configuration to the user's config directory.
func (l *ConfigLoader) SaveConfig(cfg *Config) error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get user config dir: %w", err)
	}

	appConfigDir := filepath.Join(configDir, ConfigDirName)
	if err := os.MkdirAll(appConfigDir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	configPath := filepath.Join(appConfigDir, ConfigFileName)
	
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write with 0600 permissions (User Read/Write only)
	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func GetConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, ConfigDirName, ConfigFileName), nil
}
