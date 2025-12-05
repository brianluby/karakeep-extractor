# CLI Command Contract: Setup

## Command

```bash
karakeep setup
```

## Description

Interactively prompts the user for configuration details and saves them to `config.yaml`.

## Interaction Flow

```text
$ karakeep setup

Karakeep Extractor Setup
------------------------
Enter Karakeep URL [https://karakeep.example.com]: <input>
Enter Karakeep API Token: <input>
Enter GitHub Personal Access Token (optional): <input>
Enter SQLite Database Path [./karakeep.db]: <input>

Configuration saved to /Users/user/.config/karakeep/config.yaml
Permissions set to 0600.
```

## Error States

*   **Permission Denied**: "Error: Cannot write to config directory. Please check permissions."
