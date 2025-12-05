# Quickstart: Configuration Wizard

## Interactive Setup

1.  Run the setup wizard:
    ```bash
    karakeep setup
    ```
2.  Follow the prompts to enter your URLs and Tokens.
3.  Subsequent commands (`extract`, `rank`) will now work without extra flags:
    ```bash
    karakeep extract
    karakeep rank
    ```

## Manual Configuration

You can also manually edit the file at `~/.config/karakeep/config.yaml`:

```yaml
karakeep_url: "..."
karakeep_token: "..."
```
