# Quickstart: Advanced Querying

## Query Command

Use the `query` command to filter your repository database.

### Examples

**Filter by Stars (Greater than)**
```bash
karakeep query --stars ">1000"
```

**Filter by Date (Since Jan 1, 2024)**
```bash
karakeep query --after "2024-01-01"
```

**Complex Combination**
Find Go projects with over 500 stars found this year.
```bash
karakeep query --lang Go --stars ">500" --after "2024-01-01"
```

**Filter by Forks (Range)**
```bash
karakeep query --forks "50..200"
```
