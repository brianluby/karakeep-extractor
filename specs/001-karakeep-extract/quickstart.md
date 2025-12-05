# Quickstart: Karakeep Extractor

## Prerequisites
*   Go 1.25+
*   Access to a Karakeep instance (URL + Token)

## Build
```bash
go build -o karakeep-extractor cmd/extractor/main.go
```

## Run Extraction

1.  **Set Environment Variables** (Optional but recommended):
    ```bash
    export KARAKEEP_URL="https://my-karakeep.com"
    export KARAKEEP_TOKEN="my-secret-token"
    ```

2.  **Run the command**:
    ```bash
    ./karakeep-extractor extract
    ```

3.  **Verify Output**:
    ```bash
    sqlite3 karakeep.db "SELECT * FROM extracted_repos;"
    ```

## Troubleshooting
*   **401 Unauthorized**: Check your token.
*   **Connection Refused**: Check your URL.
*   **DB Locked**: Ensure no other process is writing to `karakeep.db`.
