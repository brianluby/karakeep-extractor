# Internal Interfaces: Configuration

No new interfaces required strictly, but `ConfigLoader` logic will expose methods.

```go
// In internal/config/loader.go

// LoadConfig loads configuration with precedence: Flag > Env > File.
// flagConfig contains values parsed from CLI flags (if any).
func LoadConfig(flagConfig Config) (*Config, error)

// SaveConfig persists the configuration to the default file location.
func SaveConfig(cfg *Config) error
```
