# Developer Runbook

This document serves as the comprehensive guide for developing, testing, and releasing the `concat` project.

## 1. Project Architecture

The project follows a modular CLI architecture, splitting concerns into distinct packages:

*   `cmd/`: Entry points (`concat`, `opt`). Responsible for CLI flag parsing (Cobra).
*   `internal/app/`: High-level orchestration. Connects config, core logic, and UI.
*   `internal/core/`: Pure domain logic (File walking, Tree generation, Filtering).
*   `internal/transform/`: Stream processing logic for `opt` (Token counting, Regex replacements).
*   `internal/protocol/`: Output formatting strategies (Markdown vs XML).

## 2. Local Development

### Prerequisites
*   Go 1.21+

### Build Commands
We use standard Go tooling.

```bash
# Build concat
go build -o bin/concat cmd/concat/main.go

# Build opt
go build -o bin/opt cmd/opt/main.go
```

### Running from Source
You can test the "Power Pipe" without building binaries:

```bash
go run cmd/concat/main.go -p go | go run cmd/opt/main.go --cost
```

## 3. Testing Strategy

We employ both unit tests and end-to-end (E2E) integration tests.

### Unit Tests
Run these frequently during development.
```bash
go test ./...
```

### End-to-End (E2E) Tests
Located in `tests/e2e_test.go`. These tests compile the binaries and execute them against a real file system to verify the pipeline behaves as expected.

**Run E2E Tests:**
```bash
go test -v ./tests/...
```

## 4. Release Process

Releases are automated via [GoReleaser](https://goreleaser.com/).

### Pre-Release Checklist
1.  Run `go mod tidy`.
2.  Ensure all tests pass: `go test ./...` and `go test -v ./tests/...`.
3.  Update `README.md` if arguments or flags have changed.

### Creating a Release

1.  **Tag the version:**
    Semantic versioning is strictly enforced (vX.Y.Z).
    ```bash
    git tag -a v0.1.6 -m "feat: add XML support"
    git push origin v0.1.6
    ```

2.  **GoReleaser:**
    The CI/CD pipeline (GitHub Actions) should handle this automatically on tag push. To run locally:
    ```bash
    # Snapshot release (does not publish)
    goreleaser release --snapshot --clean
    ```

## 5. Troubleshooting

### Clipboard Issues
*   **Linux (Wayland):** `wl-copy` must be installed.
*   **Linux (X11):** `xclip` or `xsel` must be installed.
*   **Headless Environments:** The tool will detect if no clipboard is available and may fallback to stdout or error. Use `-s` / `--stdout` in scripts.

### Performance Profiling
If directory traversal seems slow:
```bash
go run cmd/concat/main.go -p go --cpuprofile cpu.prof
go tool pprof cpu.prof
```