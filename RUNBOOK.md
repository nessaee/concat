# concat & opt Developer Runbook

This runbook distills the essential workflows for developing, building, and releasing the `concat` suite.

## 1. Local Development & Testing

**Build binaries:**
```bash
# Build both tools
go build -o concat cmd/concat/main.go
go build -o opt cmd/opt/main.go
```

**Run from source (Pipeline):**
```bash
go run cmd/concat/main.go -p go | go run cmd/opt/main.go -c
```

**Run Unit Tests:**
```bash
go test ./...
```

**Run E2E Tests (Integration):**
This compiles fresh binaries and tests the pipe logic.
```bash
go test -v ./tests/...
```

## 2. Installation (Local)

To make the suite accessible globally on your machine:

```bash
# Install both tools to $GOPATH/bin
go install ./cmd/concat
go install ./cmd/opt
```

**Verify:**
```bash
concat --help
opt --help
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
    git tag -a v0.1.2 -m "Release v0.1.2"

    # 2. Push tag
    git push origin v0.1.2

    # 3. Release
    goreleaser release --clean
    ```

## 4. Troubleshooting

*   **Clipboard issues (Linux):** Ensure `wl-copy` (Wayland) or `xclip` (X11) is installed.
    *   `sudo apt install wl-clipboard` or `sudo apt install xclip`
*   **Permission Denied:** Ensure `~/.local/bin/concat` has execution permissions (`chmod +x ~/.local/bin/concat`).

## 5. Optimization & Cost Saving

To minimize token usage when feeding LLMs (like Gemini or Claude), follow these guidelines:

**1. Use `opt` (Companion Tool):**
Strip excessive whitespace and headers.
```bash
concat -p go | opt --compact --strip-headers
```

**2. Skip the Tree:**
If you are using tools like `forge` or `files-to-prompt`, the XML/header structure is enough.
```bash
# EXPENSIVE (Redundant structure)
concat -p go -t | forge

# OPTIMIZED (Structure inferred from file paths)
concat -p go | forge
```
