# Data Model: Configuration

## Configuration File (`config.yaml`)

```yaml
karakeep_url: "https://my-instance.com"
karakeep_token: "secret-token"
github_token: "ghp_..." # Optional
db_path: "/home/user/.local/share/karakeep.db"
```

## Entities

### AppConfig (Go Struct)

Updates existing `config.Config` struct to support YAML tags.

```go
type Config struct {
    KarakeepURL   string `yaml:"karakeep_url"`
    KarakeepToken string `yaml:"karakeep_token"`
    GitHubToken   string `yaml:"github_token,omitempty"`
    DBPath        string `yaml:"db_path"`
}
```
