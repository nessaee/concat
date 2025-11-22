# concat Developer Runbook

This runbook distills the essential workflows for developing, building, and releasing the `concat` tool.

## 1. Local Development & Testing

**Build binary:**
```bash
go build -o concat cmd/concat/main.go
```

**Run from source:**
```bash
go run cmd/concat/main.go -p go
```

**Run Tests:**
```bash
go test ./...
```

## 2. Installation (Local)

To make `concat` accessible globally on your machine:

```bash
# 1. Build optimized binary
CGO_ENABLED=0 go build -ldflags="-s -w" -o concat cmd/concat/main.go

# 2. Move to path (requires ~/.local/bin to be in $PATH)
mv concat ~/.local/bin/
```

**Verify:**
```bash
concat --help
```

## 3. Release Process (GoReleaser)

**Prerequisites:**
- `goreleaser` installed (`go install github.com/goreleaser/goreleaser/v2@latest`)
- `GITHUB_TOKEN` environment variable set.

**Steps:**

1.  **Snapshot Release (Test Build):**
    Builds artifacts in `dist/` without publishing.
    ```bash
    goreleaser release --snapshot --clean
    ```

2.  **Official Release:**
    Creates a tag, pushes it, and publishes binaries to GitHub Releases.
    ```bash
    # 1. Tag the version
    git tag -a v0.1.1 -m "Release v0.1.1"

    # 2. Push tag
    git push origin v0.1.1

    # 3. Release
    goreleaser release --clean
    ```

## 4. Troubleshooting

*   **Clipboard issues (Linux):** Ensure `wl-copy` (Wayland) or `xclip` (X11) is installed.
    *   `sudo apt install wl-clipboard` or `sudo apt install xclip`
*   **Permission Denied:** Ensure `~/.local/bin/concat` has execution permissions (`chmod +x ~/.local/bin/concat`).

## 5. Optimization & Cost Saving

To minimize token usage when feeding LLMs (like Gemini or Claude), follow these guidelines:

**1. Use `--compact`:**
Strip excessive whitespace.
```bash
concat -p go --compact
```

**2. Skip the Tree:**
If you are using tools like `forge` or `files-to-prompt`, the XML/header structure is enough.
```bash
# EXPENSIVE (Redundant structure)
concat -p go -t | forge

# OPTIMIZED (Structure inferred from file paths)
concat -p go | forge
```
