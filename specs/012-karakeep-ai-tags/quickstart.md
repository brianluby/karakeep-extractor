# Quickstart: Karakeep AI Tags

## 1. Extraction with Tags
Simply run the standard extract command. The system now automatically pulls tags from your bookmarks.

```bash
karakeep extract
```

**Verification**: Check the database or use the rank command to see if tags are being populated.

## 2. Filtering by Tag
Use the `--tag` flag with the `rank` command to filter repositories.

```bash
# Find all repositories tagged with "golang"
karakeep rank --tag golang
```

```bash
# Find all repositories tagged with "cli"
karakeep rank --tag cli
```

If no repositories match, the output will be empty.
