# concat

`concat` (v1.0.0) is a high-performance, cross-platform CLI tool designed to aggregate project files into a single text block. It is specifically engineered to streamline the process of gathering code context for Large Language Models (LLMs).

## Features

- **Recursive Traversal:** Efficiently walks directory trees.
- **Smart Filtering:** Respects `.gitignore` (recursively) and default ignore patterns (e.g., `node_modules`, `.git`).
- **Clipboard Integration:** Automatically copies the output to your system clipboard (supports Linux, macOS, Windows, and WSL).
- **Tree Visualization:** Generates a visual directory tree to provide structural context.
- **Single Binary:** No dependencies required; runs as a static binary.

## Installation

### From Source (Go)

If you have Go 1.21+ installed:

```bash
go install github.com/nessaee/concat/cmd/concat@latest
```

### Manual Installation (Linux)

1.  Build the project:
    ```bash
    git clone <repo_url>
    cd concat
    go build -ldflags="-s -w" -o concat cmd/concat/main.go
    ```
2.  Move to your path:
    ```bash
    sudo mv concat /usr/local/bin/
    ```

## Usage

```bash
concat -p <extension> [flags]
```

### Flags

| Flag | Description | Example |
|------|-------------|---------|
| `-p, --pattern` | **Required.** File extensions to include (no dot needed). Can be repeated. | `-p go -p md` |
| `-i, --ignore` | Glob pattern to ignore (files or directories). Can be repeated. | `-i "test/*" -i "*.log"` |
| `-t, --tree` | Include a visual directory tree at the top of the output. | `-t` |
| `-o, --output` | Write output to a file instead of the clipboard. | `-o context.txt` |
| `-h, --help` | Show help message. | |

### Examples

**1. Prepare context for a Go project:**
Captures all `.go` and `.mod` files, includes the directory structure, and copies to clipboard.
```bash
concat -p go -p mod -t
```

**2. Analyze a React/TypeScript project:**
Captures `.ts`, `.tsx`, and `.css` files, ignoring the `tests` directory.
```bash
concat -p ts -p tsx -p css -i "tests/*"
```

**3. Save Python scripts to a file:**
```bash
concat -p py -o output.txt
```

## Configuration

`concat` automatically respects your `.gitignore` files. It also has built-in defaults to ignore common noise:
- `.git`, `node_modules`, `__pycache__`, `.venv`, `target`, `dist`, `build`
- `*.log`, `*.lock`, `*.swp`, `.DS_Store`

## Development

### Build Locally
```bash
go build -o concat cmd/concat/main.go
```

### Run Tests
*Coming soon.*